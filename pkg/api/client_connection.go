package api

import (
	"context"
	"encoding/json"
	"errors"
	"path"

	"github.com/formicidae-tracker/olympus/pkg/tm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type metadated interface {
	MD() map[string]string
}

func (m *ClimateUpStream) MD() map[string]string {
	return m.Metadata
}

func (m *ClimateDownStream) MD() map[string]string {
	return m.Metadata
}

func (m *TrackingUpStream) MD() map[string]string {
	return m.Metadata
}

func (m *TrackingDownStream) MD() map[string]string {
	return m.Metadata
}

type textCarrier struct {
	m metadated
}

func (c textCarrier) Get(key string) string {
	return c.m.MD()[key]
}

func (c textCarrier) Set(key, value string) {
	c.m.MD()[key] = value
}

func (c textCarrier) Keys() []string {
	md := c.m.MD()
	res := make([]string, 0, len(md))
	for k := range md {
		res = append(res, k)
	}
	return res
}

// Stream defines a generic interface for a long-lived ping-pong stream such as
// OlympusClimate and OlympusTracking
type Stream[Up, Down any] interface {
	Send(Up) error
	Recv() (Down, error)
	CloseSend() error
	Context() context.Context
}

type connectionConfig struct {
	name        string
	tracer      trace.Tracer
	propagator  propagation.TextMapPropagator
	dialOptions []grpc.DialOption
}

type ConnectionOption interface {
	apply(*connectionConfig)
}

type withTracer string

func (o withTracer) apply(config *connectionConfig) {
	if tm.Enabled() == false {
		return
	}
	config.name = string(o)
	config.tracer = otel.GetTracerProvider().
		Tracer("github.com/formicidae-tracker/olympus/pkg/tm")
	config.propagator = otel.GetTextMapPropagator()
}

func withSpanBasename(name string) ConnectionOption {
	return withTracer(name)
}

type withDialOptions []grpc.DialOption

func (o withDialOptions) apply(config *connectionConfig) {
	config.dialOptions = append(config.dialOptions, []grpc.DialOption(o)...)
}

func WithDialOptions(opts ...grpc.DialOption) ConnectionOption {
	return withDialOptions(opts)
}

// Connection represents a connection to a long-lived ping-pong stream
// such as OlympusTracking or OlympusClimate.
type Connection[Up, Down metadated] struct {
	conn        *grpc.ClientConn
	stream      Stream[Up, Down]
	acknowledge Down
	config      connectionConfig
	links       []trace.Link
}

func (c *Connection[Up, Down]) startTrace(name string) (context.Context, trace.Span) {
	return c.config.tracer.Start(
		// explicitely use background to not endup in the grpc long-lived trace
		context.Background(),
		path.Join(c.config.name, name),
		trace.WithLinks(c.links...),
	)
}

// Established returns if the Connection is well established and can
// be used.
func (c *Connection[Up, Down]) Established() bool {
	return c.conn != nil && c.stream != nil
}

// Comfirmation returns the message of the Olympus server received on
// declaration.
func (c *Connection[Up, Down]) Confirmation() Down {
	return c.acknowledge
}

func safeEncode(v any) string {
	res, _ := json.Marshal(v)
	return string(res)
}

// Send sends a message to the Olympus implementation and returns
// either its response or an error.
func (c *Connection[Up, Down]) Send(m Up) (res Down, err error) {
	if c.stream == nil {
		return
	}

	if c.config.tracer != nil {
		ctx, span := c.startTrace("UpDownExchange")
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()

		c.config.propagator.Inject(ctx, textCarrier{m})
	}

	err = c.stream.Send(m)
	if err != nil {
		return
	}
	res, err = c.stream.Recv()
	return
}

// Close closes the connection, and report any connection error.
func (c *Connection[Up, Down]) Close() (err error) {
	var serr, cerr error
	if c.stream != nil {
		if c.config.tracer != nil {
			_, span := c.startTrace("CloseSend")
			defer func() {
				if err != nil {
					span.SetStatus(codes.Error, err.Error())
				}
				span.End()
			}()
		}
		serr = c.stream.CloseSend()
	}
	if c.conn != nil {
		cerr = c.conn.Close()
	}
	c.stream = nil
	c.conn = nil
	return errors.Join(serr, cerr)
}

// Connection results holds the result of an attemp to create a
// connection to an olympus implementation. Either one of the two
// field will be populated.
type ConnectionResult[Up, Down metadated] struct {
	Connection *Connection[Up, Down]
	Error      error
}

// connectionFactory is an helper type to instantiate either
// OlympusClimate or OlympusTracking stream.
type connectionFactory[Up, Down metadated] func(context.Context, OlympusClient) (Stream[Up, Down], Up, error)

func connect[Up, Down metadated](
	ctx context.Context,
	address string,
	factory connectionFactory[Up, Down],
	options ...ConnectionOption) (c *Connection[Up, Down], err error) {

	c = &Connection[Up, Down]{
		config: connectionConfig{
			dialOptions: DefaultDialOptions,
		},
	}

	for _, opt := range options {
		opt.apply(&c.config)
	}

	defer func() {
		if err != nil {
			c.Close()
			c = nil
		}
	}()

	c.conn, err = grpc.DialContext(ctx, address, c.config.dialOptions...)
	if err != nil {
		return
	}

	client := NewOlympusClient(c.conn)
	var decl Up
	c.stream, decl, err = factory(ctx, client)

	if err != nil {
		return
	}

	if c.config.tracer != nil {
		ctx, span := c.startTrace("Declaration")
		c.config.propagator.Inject(ctx, textCarrier{decl})
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()

		c.links = append(c.links,
			trace.Link{SpanContext: span.SpanContext()},
		)

		grpcSpanContext := trace.SpanContextFromContext(c.stream.Context())
		if grpcSpanContext.IsValid() {
			c.links = append(c.links, trace.Link{SpanContext: grpcSpanContext})
		}
	}

	err = c.stream.Send(decl)
	if err != nil {
		return
	}

	c.acknowledge, err = c.stream.Recv()

	return
}

// Connect connects asynchronously to an Olympus implementation and
// returns the associated ConnectionResult.
func Connect[Up, Down metadated](
	ctx context.Context,
	address string,
	factory connectionFactory[Up, Down],
	options ...ConnectionOption) <-chan ConnectionResult[Up, Down] {

	res := make(chan ConnectionResult[Up, Down])

	go func() {
		defer close(res)

		var connResult ConnectionResult[Up, Down]

		connResult.Connection, connResult.Error = connect(
			ctx, address, factory, options...)

		res <- connResult
	}()
	return res
}

func climateConnector(declaration *ClimateDeclaration) connectionFactory[*ClimateUpStream, *ClimateDownStream] {

	return func(ctx context.Context,
		client OlympusClient) (Stream[*ClimateUpStream, *ClimateDownStream], *ClimateUpStream, error) {

		stream, err := client.Climate(ctx, DefaultCallOptions...)

		return stream, &ClimateUpStream{Declaration: declaration}, err

	}

}

func trackingConnector(declaration *TrackingDeclaration) connectionFactory[*TrackingUpStream, *TrackingDownStream] {

	return func(ctx context.Context,
		client OlympusClient) (Stream[*TrackingUpStream, *TrackingDownStream], *TrackingUpStream, error) {

		stream, err := client.Tracking(ctx, DefaultCallOptions...)

		return stream, &TrackingUpStream{Declaration: declaration}, err
	}
}

// ConnectClimate asynchronously connect and call an OlympusClimate
// stream on a given address, with the given declaration.
func ConnectClimate(
	ctx context.Context,
	address string,
	declaration *ClimateDeclaration,
	options ...ConnectionOption) <-chan ConnectionResult[*ClimateUpStream, *ClimateDownStream] {

	return Connect(ctx,
		address,
		climateConnector(declaration),
		append(options, withSpanBasename("fort.olympus.Olympus/Climate"))...)
}

// ConnectTracking asynchronously connect and call an OlympusTracking
// stream on a given address, with the given declaration.
func ConnectTracking(
	ctx context.Context,
	address string,
	declaration *TrackingDeclaration,
	options ...ConnectionOption) <-chan ConnectionResult[*TrackingUpStream, *TrackingDownStream] {

	return Connect(ctx,
		address,
		trackingConnector(declaration),
		append(options, withSpanBasename("fort.olympus.Olympus/Tracking"))...)
}
