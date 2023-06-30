package api

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

// RequestResult represents the result of a Request to a ClientTask.
type RequestResult[Down any] struct {
	Message *Down
	Error   error
}

// Request represent a request and its expected Result to a ClientTask
type Request[Up, Down any] struct {
	Message  *Up
	response chan RequestResult[Down]
}

// NewRequest prepares a new Request for a ClientTask
func NewRequest[Up, Down any](m *Up) Request[Up, Down] {
	return Request[Up, Down]{
		Message:  m,
		response: make(chan RequestResult[Down], 1),
	}
}

// Response returns the asynchronous RequestResult of a Request.
func (r Request[Up, Down]) Response() <-chan RequestResult[Down] {
	return r.response
}

// ConfirmationResult returns the result of a new connection to a
// Olympus stream such as OlympusClimate or OlympusTracking.
type ConfirmationResult[Down any] struct {
	Confirmation *Down
	Error        error
}

// ClientTask is a task that can be used to manage a long-lived
// connection to an Olympus stream. It will attempt to reconnect
// automatically to the Olympus server. It will sends any connection
// attempt result on its Confirmations() channel. Request() can be
// used to perform an asynchronous request with this task. Run()
// should be called to perform the task loop.
type ClientTask[Up, Down any] struct {
	ctx           context.Context
	cancel        context.CancelFunc
	connect       func() <-chan ConnectionResult[Up, Down]
	inbound       chan Request[Up, Down]
	confirmations chan ConfirmationResult[Down]
}

// Run runs the ClientTask communication loops until Close() is called.
func (c *ClientTask[Up, Down]) Run() (err error) {
	var connection *Connection[Up, Down] = nil
	defer func() {
		if connection == nil {
			return
		}
		cerr := connection.Close()
		if err == nil {
			err = cerr
		} else if cerr != nil {
			err = fmt.Errorf("multiple errors: %s", []error{err, cerr})
		}
	}()
	defer close(c.confirmations)

	newConnection := c.connect()

	for {
		if newConnection == nil {
			c.sleepWithJitter(2*time.Second, 0.1)
			newConnection = c.connect()
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
			} else {
				connection = res.Connection
				confirmation.Confirmation = connection.acknowledge
			}

			select {
			case c.confirmations <- confirmation:
			default:
			}

		case req, ok := <-c.inbound:
			if ok == false {
				return
			}
			if connection == nil {
				continue
			}
			var res RequestResult[Down]
			res.Message, res.Error = connection.Send(req.Message)
			req.response <- res
		}

	}
}

func (c *ClientTask[Up, Down]) sleepWithJitter(d time.Duration, jitter float64) {
	toSleep := (1.0 + (2.0*rand.Float64()-1.0)*jitter) * float64(d.Nanoseconds())
	time.Sleep(time.Duration(toSleep))
}

// Stop stops the Run loop.
func (c *ClientTask[Up, Down]) Stop() {
	c.cancel()
}

// Request performs an asynchronous Request on the ClientTask
func (c *ClientTask[Up, Down]) Request(m *Up) <-chan RequestResult[Down] {
	req := NewRequest[Up, Down](m)
	c.inbound <- req
	return req.Response()
}

// Confirmations returns all Connection attempt result to the Olympus
// stream.
func (c *ClientTask[Up, Down]) Confirmations() <-chan ConfirmationResult[Down] {
	return c.confirmations
}

func newClientTask[Up, Down any](
	ctx context.Context,
	address string,
	factory ConnectionFactory[Up, Down],
	options ...grpc.DialOption) *ClientTask[Up, Down] {
	ctx, cancel := context.WithCancel(ctx)
	return &ClientTask[Up, Down]{
		ctx:    ctx,
		cancel: cancel,
		connect: func() <-chan ConnectionResult[Up, Down] {
			return Connect(ctx, address, factory, options...)
		},
		inbound:       make(chan Request[Up, Down], 10),
		confirmations: make(chan ConfirmationResult[Down], 2),
	}
}

// NewClimateTask returns a ClimateTask that connect and call an
// OlympusClimate stream.
func NewClimateTask(
	ctx context.Context,
	address string,
	declaration *ClimateDeclaration,
	options ...grpc.DialOption) *ClientTask[ClimateUpStream, ClimateDownStream] {
	return newClientTask(
		ctx, address,
		climateConnector(declaration),
		options...)
}

// NewTrackingTask returns a ClimateTask that connect and call an
// OlympusTracking stream.
func NewTrackingTask(
	ctx context.Context,
	address string,
	declaration *TrackingDeclaration,
	options ...grpc.DialOption) *ClientTask[TrackingUpStream, TrackingDownStream] {
	return newClientTask(
		ctx, address,
		trackingConnector(declaration),
		options...)
}
