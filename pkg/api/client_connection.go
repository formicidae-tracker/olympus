package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Stream defines a generic interface for a long-lived ping-pong stream such as
// OlympusClimate and OlympusTracking
type Stream[Up any, Down any] interface {
	Send(*Up) error
	Recv() (*Down, error)
	CloseSend() error
}

// Connection represents a connection to a long-lived ping-pong stream
// such as OlympusTracking or OlympusClimate.
type Connection[Up any, Down any] struct {
	conn        *grpc.ClientConn
	stream      Stream[Up, Down]
	acknowledge *Down
}

// Established returns if the Connection is well established and can
// be used.
func (c *Connection[Up, Down]) Established() bool {
	return c.conn != nil && c.stream != nil
}

// Comfirmation returns the message of the Olympus server received on
// declaration.
func (c *Connection[Up, Down]) Confirmation() *Down {
	return c.acknowledge
}

// Send sends a message to the Olympus implementation and returns
// either its response or an error.
func (c *Connection[Up, Down]) Send(m *Up) (*Down, error) {
	if c.stream == nil {
		return nil, nil
	}
	err := c.stream.Send(m)
	if err != nil {
		return nil, err
	}
	return c.stream.Recv()
}

// Close closes the connection, and report any connection error.
func (c *Connection[Up, Down]) Close() (err error) {
	var errs []error = nil
	defer func() {
		if len(errs) == 0 {
			err = nil
		} else {
			err = fmt.Errorf("multiple errors: %s", errs)
		}
	}()

	if c.stream != nil {
		cerr := c.stream.CloseSend()
		if cerr != nil {
			errs = append(errs, cerr)
		}
	}
	c.stream = nil
	if c.conn == nil {
		return
	}
	cerr := c.conn.Close()
	if cerr != nil {
		errs = append(errs, cerr)
	}

	return
}

// Connection results holds the result of an attemp to create a
// connection to an olympus implementation. Either one of the two
// field will be populated.
type ConnectionResult[Up, Down any] struct {
	Connection *Connection[Up, Down]
	Error      error
}

// ConnectionFactory is an helper type to instantiate either
// OlympusClimate or OlympusTracking stream.
type ConnectionFactory[Up, Down any] func(context.Context, OlympusClient) (Stream[Up, Down], *Down, error)

func connect[Up, Down any](
	ctx context.Context,
	address string,
	factory ConnectionFactory[Up, Down],
	options ...grpc.DialOption) (c *Connection[Up, Down], err error) {
	c = &Connection[Up, Down]{}
	defer func() {
		if err != nil {
			c.Close()
		}
		c = nil
	}()

	options = append(DefaultDialOptions, options...)
	c.conn, err = grpc.DialContext(ctx, address, options...)
	if err != nil {
		return
	}
	client := NewOlympusClient(c.conn)
	c.stream, c.acknowledge, err = factory(ctx, client)

	return
}

// Connect connects asynchronously to an Olympus implementation and
// returns the associated ConnectionResult.
func Connect[Up, Down any](
	ctx context.Context,
	address string,
	factory ConnectionFactory[Up, Down],
	options ...grpc.DialOption) <-chan ConnectionResult[Up, Down] {
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

func climateConnector(declaration *ClimateDeclaration) ConnectionFactory[ClimateUpStream, ClimateDownStream] {

	return func(ctx context.Context, client OlympusClient) (
		Stream[ClimateUpStream, ClimateDownStream], *ClimateDownStream, error) {

		stream, err := client.Climate(ctx, DefaultCallOptions...)
		if err != nil {
			return stream, nil, err
		}
		err = stream.Send(&ClimateUpStream{
			Declaration: declaration,
		})
		if err != nil {
			return stream, nil, err
		}
		acknowledge, err := stream.Recv()

		return stream, acknowledge, err
	}

}

func trackingConnector(declaration *TrackingDeclaration) ConnectionFactory[TrackingUpStream, TrackingDownStream] {

	return func(ctx context.Context, client OlympusClient) (
		Stream[TrackingUpStream, TrackingDownStream], *TrackingDownStream, error) {

		stream, err := client.Tracking(ctx, DefaultCallOptions...)
		if err != nil {
			return stream, nil, err
		}
		err = stream.Send(&TrackingUpStream{
			Declaration: declaration,
		})
		if err != nil {
			return stream, nil, err
		}
		acknowledge, err := stream.Recv()

		return stream, acknowledge, err
	}

}

// ConnectClimate asynchronously connect and call an OlympusClimate
// stream on a given address, with the given declaration.
func ConnectClimate(
	ctx context.Context,
	address string,
	declaration *ClimateDeclaration,
	options ...grpc.DialOption) <-chan ConnectionResult[ClimateUpStream, ClimateDownStream] {

	return Connect(ctx, address, climateConnector(declaration), options...)
}

// ConnectTracking asynchronously connect and call an OlympusTracking
// stream on a given address, with the given declaration.
func ConnectTracking(
	ctx context.Context,
	address string,
	declaration *TrackingDeclaration,
	options ...grpc.DialOption) <-chan ConnectionResult[TrackingUpStream, TrackingDownStream] {

	return Connect(ctx, address, trackingConnector(declaration), options...)
}
