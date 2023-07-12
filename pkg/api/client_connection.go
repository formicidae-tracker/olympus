package api

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/tm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// A metadated is a type that holds a map[string]string. It is used to
// inject trace.SpanContext in Up/Down stream message.
type metadated interface {
	md(create bool) map[string]string
}

func (m *ClimateUpStream) md(create bool) map[string]string {
	if create == true && m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
	return m.Metadata
}

func (m *ClimateDownStream) md(create bool) map[string]string {
	if create == true && m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
	return m.Metadata
}

func (m *TrackingUpStream) md(create bool) map[string]string {
	if create == true && m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
	return m.Metadata
}

func (m *TrackingDownStream) md(create bool) map[string]string {
	if create == true && m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
	return m.Metadata
}

// A textCarrier implements trace.TextCarrier for a metadated.
type textCarrier struct {
	m metadated
}

func (c textCarrier) Get(key string) string {
	md := c.m.md(false)
	if md == nil {
		return ""
	}
	return md[key]
}

func (c textCarrier) Set(key, value string) {
	c.m.md(true)[key] = value
}

func (c textCarrier) Keys() []string {
	md := c.m.md(false)
	res := make([]string, 0, len(md))
	for k := range md {
		res = append(res, k)
	}
	return res
}

// A stream defines a generic interface for a long-lived ping-pong
// stream such as [OlympusClient.Climate] and
// [OlympusClient.Tracking].
type stream[Up, Down any] interface {
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
	delay       time.Duration
}

// A ConnectionOption represents an optional parameter to [Connect] or
// for a [ClientTask].
type ConnectionOption interface {
	apply(*connectionConfig)
}

type connectionOptionFunc func(*connectionConfig)

func (f connectionOptionFunc) apply(c *connectionConfig) {
	f(c)
}

// Adds a trace base name for stream operation (Declaration,
// UpDownExchange, Close)
func withSpanBasename(name string) ConnectionOption {
	return connectionOptionFunc(func(config *connectionConfig) {
		if tm.Enabled() == false {
			return
		}
		config.name = string(name)
		config.tracer = otel.GetTracerProvider().
			Tracer("github.com/formicidae-tracker/olympus/pkg/tm")
		config.propagator = otel.GetTextMapPropagator()
	})
}

func withDelay(delay time.Duration) ConnectionOption {
	return connectionOptionFunc(func(opts *connectionConfig) {
		opts.delay = delay
	})
}

type withDialOptions []grpc.DialOption

func (o withDialOptions) apply(config *connectionConfig) {
	config.dialOptions = append(config.dialOptions, []grpc.DialOption(o)...)
}

// WithDialOptions adds the opts grpc.DialOption to the connection
// config.
func WithDialOptions(opts ...grpc.DialOption) ConnectionOption {
	return connectionOptionFunc(func(config *connectionConfig) {
		config.dialOptions = append(config.dialOptions, opts...)
	})
}

// A Connection represents a connection to a long-lived ping-pong
// stream provided by [OlympusClient].
type Connection[Up, Down metadated] struct {
	conn        *grpc.ClientConn
	stream      stream[Up, Down]
	acknowledge Down
	config      connectionConfig
	links       []trace.Link
}

// starts a trace for the connection.
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

func endSpanWithError(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

// Send sends a message to the Olympus implementation and returns
// either its response or an error.
func (c *Connection[Up, Down]) Send(m Up) (res Down, err error) {
	if c.stream == nil {
		return
	}

	if c.config.tracer != nil {
		ctx, span := c.startTrace("UpDownExchange")
		defer func() { endSpanWithError(span, err) }()

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
	var errs []error
	if c.stream != nil {
		if c.config.tracer != nil {
			_, span := c.startTrace("CloseSend")
			defer func() { endSpanWithError(span, err) }()
		}
		errs = append(errs, c.stream.CloseSend())
	}
	if c.conn != nil {
		errs = append(errs, c.conn.Close())
	}
	c.stream = nil
	c.conn = nil
	return errors.Join(errs...)
}

// ConnectionResult holds the result of an attemp to create a
// connection to an olympus implementation through [Connect]. Either
// one of the two field will be populated.
type ConnectionResult[Up, Down metadated] struct {
	Connection *Connection[Up, Down]
	Error      error
}

func newConnectionOptions(options ...ConnectionOption) connectionConfig {
	res := connectionConfig{
		dialOptions: append(DefaultDialOptions, grpc.WithBlock()),
	}

	for _, option := range options {
		option.apply(&res)
	}

	return res
}

// connectionFactory is an helper type to instantiate either
// OlympusClimate or OlympusTracking stream.
type connectionFactory[Up, Down metadated] func(context.Context, OlympusClient) (stream[Up, Down], Up, error)

// connectSync connects blockly connectSync to an [OlympusClient] stream.
func connectSync[Up, Down metadated](
	ctx context.Context,
	address string,
	factory connectionFactory[Up, Down],
	options ...ConnectionOption) (c *Connection[Up, Down], err error) {

	c = &Connection[Up, Down]{
		config: newConnectionOptions(options...),
	}

	defer func() {
		if err != nil {
			c.Close()
			c = nil
		}
	}()

	// links with connect context
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() == true {
		c.links = append(c.links, trace.Link{
			SpanContext: sc,
		})
	}

	time.Sleep(c.config.delay)

	c.conn, err = grpc.DialContext(ctx, address,
		c.config.dialOptions...)
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
		defer func() { endSpanWithError(span, err) }()

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

// connect connects asynchronously to an Olympus implementation and
// returns the associated ConnectionResult.
func connect[Up, Down metadated](
	ctx context.Context,
	address string,
	factory connectionFactory[Up, Down],
	options ...ConnectionOption) <-chan ConnectionResult[Up, Down] {

	res := make(chan ConnectionResult[Up, Down])

	go func() {
		defer close(res)

		var connResult ConnectionResult[Up, Down]

		connResult.Connection, connResult.Error = connectSync(
			ctx, address, factory, options...)

		res <- connResult
	}()
	return res
}

func climateConnector(declaration *ClimateDeclaration) connectionFactory[*ClimateUpStream, *ClimateDownStream] {

	return func(ctx context.Context,
		client OlympusClient) (stream[*ClimateUpStream, *ClimateDownStream], *ClimateUpStream, error) {

		stream, err := client.Climate(ctx, DefaultCallOptions...)

		return stream, &ClimateUpStream{Declaration: declaration}, err

	}

}

func trackingConnector(declaration *TrackingDeclaration) connectionFactory[*TrackingUpStream, *TrackingDownStream] {

	return func(ctx context.Context,
		client OlympusClient) (stream[*TrackingUpStream, *TrackingDownStream], *TrackingUpStream, error) {

		stream, err := client.Tracking(ctx, DefaultCallOptions...)

		return stream, &TrackingUpStream{Declaration: declaration}, err
	}
}

// ConnectClimate asynchronously connect via an [OlympusClient] and
// starts a Climate stream on a given address, with the given
// declaration and options.
func ConnectClimate(
	ctx context.Context,
	address string,
	declaration *ClimateDeclaration,
	options ...ConnectionOption) <-chan ConnectionResult[*ClimateUpStream, *ClimateDownStream] {

	return connect(ctx,
		address,
		climateConnector(declaration),
		append(options, withSpanBasename("fort.olympus.Olympus/Climate"))...)
}

// ConnectTracking asynchronously connect via an [OlympusClient] and
// starts a Tracking stream on a given address, with the given
// declaration and options.
func ConnectTracking(
	ctx context.Context,
	address string,
	declaration *TrackingDeclaration,
	options ...ConnectionOption) <-chan ConnectionResult[*TrackingUpStream, *TrackingDownStream] {

	return connect(ctx,
		address,
		trackingConnector(declaration),
		append(options, withSpanBasename("fort.olympus.Olympus/Tracking"))...)
}
