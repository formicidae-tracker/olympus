// Code generated by go generate DO NOT EDIT.

package olympuspb

import (
	"context"
	"log"

	"github.com/barkimedes/go-deepcopy"
	grpc "google.golang.org/grpc"
)

// ZoneConnection holds a connection to an olympus Zone
// bi-directional string.
type ZoneConnection struct {
	conn        *grpc.ClientConn
	stream      Olympus_ZoneClient
	acknowledge *ZoneDownStream
	log         *log.Logger
}

// Creates an new unconnected ZoneConnection.
func NewZoneConnectionWithLogger(logger *log.Logger) *ZoneConnection {
	return &ZoneConnection{
		log: logger,
	}
}

// Established returns true if connection is established.
func (c *ZoneConnection) Established() bool {
	return c.conn != nil && c.stream != nil
}

// Send sends a ZoneUpStream message and gets it ZoneDownStream
// response (typically acknowledge).
func (c *ZoneConnection) Send(m *ZoneUpStream) (*ZoneDownStream, error) {
	if c.stream == nil {
		return nil, nil
	}
	err := c.stream.Send(m)
	if err != nil {
		return nil, err
	}
	return c.stream.Recv()
}

// CloseStream close only the bi-directional string, but keeps the tcp
// connection alive.
func (c *ZoneConnection) CloseStream() {
	if c.stream != nil {
		err := c.stream.CloseSend()
		if err != nil && c.log != nil {
			c.log.Printf("gRPC CloseSend() failure: %s", err)
		}
	}
	c.stream = nil
}

// CloseAndLogErrors() close completely the ZoneConnection, avoiding
// any leaking.
func (c *ZoneConnection) CloseAndLogErrors() {
	c.CloseStream()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil && c.log != nil {
			c.log.Printf("gRPC Close() failure: %s", err)
		}
	}
	c.conn = nil
}

// Connect connects, a possibly half-connected ZoneConnection,
// and return a new connection.
func (c *ZoneConnection) Connect(address string, declaration *ZoneDeclaration, opts ...grpc.DialOption) (res *ZoneConnection, err error) {
	defer func() {
		if err == nil {
			return
		}
		res.CloseAndLogErrors()
	}()
	res.log = c.log
	if c.conn == nil {
		dialOptions := append(DefaultDialOptions, opts...)
		if res.log != nil {
			res.log.Printf("Dialing '%s'", address)
		}
		res.conn, err = grpc.Dial(address, dialOptions...)
		if err != nil {
			return
		}
	} else {
		res.conn = c.conn
	}
	client := NewOlympusClient(res.conn)

	res.stream, err = client.Zone(context.Background(), DefaultCallOptions...)
	if err != nil {
		return
	}
	res.acknowledge, err = res.Send(&ZoneUpStream{
		Declaration: declaration,
	})
	return
}

// ConnectAsync Connects a posibly half-connected ZoneConnection
// asynchronously.
func (c *ZoneConnection) ConnectAsync(address string,
	declaration *ZoneDeclaration,
	opts ...grpc.DialOption) (<-chan *ZoneConnection, <-chan error) {

	errors := make(chan error)
	connections := make(chan *ZoneConnection)
	declaration = deepcopy.MustAnything(declaration).(*ZoneDeclaration)

	go func() {
		conn, err := c.Connect(address, declaration, opts...)
		if err != nil {
			select {
			case errors <- err:
			default:
				if c.log != nil {
					c.log.Printf("gRPC connection failed after shutdown: %s", err)
				}
			}
		} else {
			select {
			case connections <- conn:
			default:
				if c.log != nil {
					c.log.Printf("gRPC connection established after shutdown. Closing it")
				}
				conn.CloseAndLogErrors()
			}
		}
	}()

	return connections, errors
}

// TrackingConnection holds a connection to an olympus Tracking
// bi-directional string.
type TrackingConnection struct {
	conn        *grpc.ClientConn
	stream      Olympus_TrackingClient
	acknowledge *TrackingDownStream
	log         *log.Logger
}

// Creates an new unconnected TrackingConnection.
func NewTrackingConnectionWithLogger(logger *log.Logger) *TrackingConnection {
	return &TrackingConnection{
		log: logger,
	}
}

// Established returns true if connection is established.
func (c *TrackingConnection) Established() bool {
	return c.conn != nil && c.stream != nil
}

// Send sends a TrackingUpStream message and gets it TrackingDownStream
// response (typically acknowledge).
func (c *TrackingConnection) Send(m *TrackingUpStream) (*TrackingDownStream, error) {
	if c.stream == nil {
		return nil, nil
	}
	err := c.stream.Send(m)
	if err != nil {
		return nil, err
	}
	return c.stream.Recv()
}

// CloseStream close only the bi-directional string, but keeps the tcp
// connection alive.
func (c *TrackingConnection) CloseStream() {
	if c.stream != nil {
		err := c.stream.CloseSend()
		if err != nil && c.log != nil {
			c.log.Printf("gRPC CloseSend() failure: %s", err)
		}
	}
	c.stream = nil
}

// CloseAndLogErrors() close completely the TrackingConnection, avoiding
// any leaking.
func (c *TrackingConnection) CloseAndLogErrors() {
	c.CloseStream()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil && c.log != nil {
			c.log.Printf("gRPC Close() failure: %s", err)
		}
	}
	c.conn = nil
}

// Connect connects, a possibly half-connected TrackingConnection,
// and return a new connection.
func (c *TrackingConnection) Connect(address string, declaration *TrackingDeclaration, opts ...grpc.DialOption) (res *TrackingConnection, err error) {
	defer func() {
		if err == nil {
			return
		}
		res.CloseAndLogErrors()
	}()
	res.log = c.log
	if c.conn == nil {
		dialOptions := append(DefaultDialOptions, opts...)
		if res.log != nil {
			res.log.Printf("Dialing '%s'", address)
		}
		res.conn, err = grpc.Dial(address, dialOptions...)
		if err != nil {
			return
		}
	} else {
		res.conn = c.conn
	}
	client := NewOlympusClient(res.conn)

	res.stream, err = client.Tracking(context.Background(), DefaultCallOptions...)
	if err != nil {
		return
	}
	res.acknowledge, err = res.Send(&TrackingUpStream{
		Declaration: declaration,
	})
	return
}

// ConnectAsync Connects a posibly half-connected TrackingConnection
// asynchronously.
func (c *TrackingConnection) ConnectAsync(address string,
	declaration *TrackingDeclaration,
	opts ...grpc.DialOption) (<-chan *TrackingConnection, <-chan error) {

	errors := make(chan error)
	connections := make(chan *TrackingConnection)
	declaration = deepcopy.MustAnything(declaration).(*TrackingDeclaration)

	go func() {
		conn, err := c.Connect(address, declaration, opts...)
		if err != nil {
			select {
			case errors <- err:
			default:
				if c.log != nil {
					c.log.Printf("gRPC connection failed after shutdown: %s", err)
				}
			}
		} else {
			select {
			case connections <- conn:
			default:
				if c.log != nil {
					c.log.Printf("gRPC connection established after shutdown. Closing it")
				}
				conn.CloseAndLogErrors()
			}
		}
	}()

	return connections, errors
}
