package api

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Error that could be returned within a [RequestResult].
type Error string

// Implements error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	// Request wasn't sent as the task is ended
	ErrorTaskEnded = Error("task ended")

	// Request wasn't sent as there was no active connection
	ErrorNoActiveConnection = Error("no active connection")

	// Request wasn't sent as the request buffer is full
	ErrorFullBuffer = Error("full request buffer")
)

// A RequestResult represents the result of a [ClientTask.Request] or
// [ClientTask.MayRequest].
type RequestResult[Down any] struct {
	// The response message from the server
	Message Down

	// The potential error. Either a gRPC error, or an Error
	Error error
}

// request represent a request and its expected Result to a ClientTask
type request[Up, Down any] struct {
	Message            Up
	response           chan RequestResult[Down]
	disposeIfNotActive bool
}

// newRequest prepares a new Request for a ClientTask
func newRequest[Up, Down any](m Up) request[Up, Down] {
	return request[Up, Down]{
		Message:  m,
		response: make(chan RequestResult[Down], 1),
	}
}

func (r request[Up, Down]) respond(m Down, err error) {
	r.response <- RequestResult[Down]{m, err}
	close(r.response)
}

// ConfirmationResult returns the result of a new connection to an
// [OlympusClient] through a [ClientTask].
type ConfirmationResult[Down any] struct {
	Confirmation Down
	Error        error
}

// A ClientTask is a task that can be used to manage a long-lived
// connection to an Olympus stream. It will attempt to reconnect
// automatically to the Olympus server via an [OlympusClient]. It will
// sends any connection attempt result on its
// [ClientTask.Confirmations] channel. [ClientTask.Request] and
// [ClientTask.MayRequest] can be called to perform an asynchronous
// request with this task. [ClientTask.Run] should be explicitely
// called to perform the task loop in the background.
type ClientTask[Up, Down metadated] struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	fatal               context.CancelCauseFunc
	connect             func(time.Duration) <-chan ConnectionResult[Up, Down]
	inbound             chan request[Up, Down]
	confirmations       chan ConfirmationResult[Down]
	connection          *Connection[Up, Down]
	buffer              []request[Up, Down]
	baseDelay, maxDelay time.Duration
}

// Run runs the [ClientTask] communication loop until either
// [ClientTask.Fatal] is called or the provided context is canceled
// (graceful exit).
func (c *ClientTask[Up, Down]) Run() (err error) {
	defer func() {
		if c.connection == nil {
			return
		}
		err = errors.Join(err, c.connection.Close())
	}()
	defer close(c.confirmations)

	newConnection := c.connect(0)

	defer func() {
		var Nil Down
		for _, req := range c.buffer {
			req.respond(Nil, ErrorTaskEnded)
		}
	}()

	i := 0

	c.resetBuffer()
	for {
		if c.connection == nil && newConnection == nil {
			newConnection = c.connect(c.delay(i))
		}

		select {
		case <-c.ctx.Done():
			return
		case res, ok := <-newConnection:
			if ok == false {
				newConnection = nil
				continue
			}
			var confirmation ConfirmationResult[Down]

			if res.Error != nil {
				confirmation.Error = res.Error
				i += 1
			} else {
				c.connection = res.Connection
				confirmation.Confirmation = c.connection.acknowledge
				hasError := false
				for i, req := range c.buffer {
					if c.handle(req) != nil {
						hasError = true
						c.buffer = c.buffer[i:]
						break
					}
				}

				if hasError == false {
					c.resetBuffer()
					i = 0
				}
			}

			select {
			case c.confirmations <- confirmation:
			default:
			}

		case req, ok := <-c.inbound:
			if ok == false {
				return
			}
			if c.connection == nil {
				c.pushOrDiscard(req)
			} else {
				c.handle(req)
			}
		}
	}
}

func (c *ClientTask[Up, Down]) delay(iteration int) time.Duration {
	delay := time.Duration(math.Pow(1.2, float64(iteration)) * float64(c.baseDelay))
	if delay > c.maxDelay {
		delay = c.maxDelay
	}

	jitter := 1.0 + 0.1*(rand.Float64()*2.0-1.0)
	return time.Duration(jitter * float64(delay))
}

var maxBufferSize = 100

func (c *ClientTask[Up, Down]) resetBuffer() {
	c.buffer = make([]request[Up, Down], 0, maxBufferSize)
}

func (c *ClientTask[Up, Down]) pushOrDiscard(req request[Up, Down]) {
	if req.disposeIfNotActive == true {
		var Nil Down
		req.respond(Nil, ErrorNoActiveConnection)
		return
	}

	for len(c.buffer) >= maxBufferSize {
		discarded := c.buffer[0]
		var Nil Down
		discarded.respond(Nil, ErrorFullBuffer)
		c.buffer = c.buffer[1:]
	}

	c.buffer = append(c.buffer, req)
}

func (c *ClientTask[Up, Down]) handle(req request[Up, Down]) error {
	if c.connection == nil {
		var Nil Down
		req.respond(Nil, ErrorNoActiveConnection)
		return ErrorNoActiveConnection
	}
	res, err := c.connection.Send(req.Message)
	if err != nil {
		err = errors.Join(err, c.connection.Close())
		c.connection = nil
	}
	req.respond(res, err)
	return err
}

// Fatal stops the [ClientTask.Run] loop with an error that will be
// propagated to the server.  If err is nil, the task will be
// gracefully terminated instead, as if the task's provided context
// was cancelled.
func (c *ClientTask[Up, Down]) Fatal(err error) {
	if err == nil {
		c.cancel()
	} else {
		c.fatal(err)
	}
}

// stop is a short hand for Fatal(nil) (graceful exit).
func (c *ClientTask[Up, Down]) stop() {
	c.Fatal(nil)
}

// Request performs an asynchronous request on the [ClientTask]. If the
// underlying connection is not active yet, it will buffer the request
// until a connection is active or the buffer is filled up. In the
// later request are canceled in FIFO order.
func (c *ClientTask[Up, Down]) Request(m Up) <-chan RequestResult[Down] {
	req := newRequest[Up, Down](m)
	go func() { c.inbound <- req }()
	return req.response
}

// MayRequest performs an asynchronous request on the ClientTask, like
// [ClientTask.Request], but discard immediatly the request if there
// is no active connection to the server.
func (c *ClientTask[Up, Down]) MayRequest(m Up) <-chan RequestResult[Down] {
	req := newRequest[Up, Down](m)
	req.disposeIfNotActive = true
	go func() { c.inbound <- req }()
	return req.response
}

// Confirmations returns all connection attempt result to the Olympus
// stream server.
func (c *ClientTask[Up, Down]) Confirmations() <-chan ConfirmationResult[Down] {
	return c.confirmations
}

func newClientTask[Up, Down metadated](
	ctx context.Context,
	address string,
	factory connectionFactory[Up, Down],
	options ...ConnectionOption) *ClientTask[Up, Down] {

	cancelable, cancel := context.WithCancel(ctx)

	connectionCtx, fatal := context.WithCancelCause(context.Background())
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() == true {
		connectionCtx = trace.ContextWithSpanContext(connectionCtx, sc)
	}

	return &ClientTask[Up, Down]{
		ctx:    cancelable,
		cancel: cancel,
		fatal:  fatal,
		connect: func(delay time.Duration) <-chan ConnectionResult[Up, Down] {
			actualOptions := append(options, withDelay(delay))
			return connect(connectionCtx, address, factory, actualOptions...)
		},
		baseDelay:     500 * time.Millisecond,
		maxDelay:      30 * time.Second,
		inbound:       make(chan request[Up, Down], 10),
		confirmations: make(chan ConfirmationResult[Down], 2),
	}
}

// NewClimateTask returns a [ClientTask] that connect and call an
// [OlympusClient] Climate stream.
func NewClimateTask(
	ctx context.Context,
	address string,
	declaration *ClimateDeclaration,
	options ...ConnectionOption) *ClientTask[*ClimateUpStream, *ClimateDownStream] {

	return newClientTask(
		ctx, address,
		climateConnector(declaration),
		append(options, withSpanBasename("fort.olympus.Olympus/Climate"))...)
}

// NewTrackingTask returns a [ClimateTask] that connect and call an
// [OlympusClient] Tracking stream.
func NewTrackingTask(
	ctx context.Context,
	address string,
	declaration *TrackingDeclaration,
	options ...ConnectionOption) *ClientTask[*TrackingUpStream, *TrackingDownStream] {

	return newClientTask(
		ctx, address,
		trackingConnector(declaration),
		append(options, withSpanBasename("fort.olympus.Olympus/Tracking"))...)
}
