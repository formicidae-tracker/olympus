//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("could not generate client utils: %s", err)
	}
}

type clientUtilTemplateData struct {
	Clients []string
}

var client_util_templates *template.Template = template.Must(template.New("").Parse(`// Code generated by go generate DO NOT EDIT.

package olympuspb

import (
	"context"
	"log"

	"github.com/barkimedes/go-deepcopy"
	grpc "google.golang.org/grpc"
)

{{ range .Clients }}

// {{. -}} Connection holds a connection to an olympus {{.}}
// bi-directional string.
type {{. -}} Connection struct {
	conn        *grpc.ClientConn
	stream      Olympus_ {{- . -}} Client
	acknowledge * {{- . -}} DownStream
	log         *log.Logger
}

// Creates an new unconnected {{ . -}} Connection.
func New {{- . -}} ConnectionWithLogger(logger *log.Logger) * {{- . -}} Connection {
	return & {{- . -}} Connection{
		log: logger,
	}
}

// Established returns true if connection is established.
func (c * {{- . -}} Connection) Established() bool {
	return c.conn != nil && c.stream != nil
}

// Confirmation returns the {{ . -}} DownStream acknowledgement for the
// olympus server declaration. It can be empty.
func (c * {{- . -}} Connection) Confirmation() * {{- . -}} DownStream {
	return c.acknowledge
}

// Send sends a {{. -}} UpStream message and gets it {{. -}} DownStream
// response (typically acknowledge).
func (c * {{- . -}} Connection) Send(m * {{- . -}} UpStream) (* {{- . -}} DownStream, error) {
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
func ( c * {{- . -}} Connection) CloseStream() {
	if c.stream != nil {
		err := c.stream.CloseSend()
		if err != nil && c.log != nil {
			c.log.Printf("gRPC CloseSend() failure: %s", err)
		}
	}
	c.stream = nil
	c.acknowledge = nil
}

// CloseAndLogErrors() close completely the {{ . -}} Connection, avoiding
// any leaking.
func (c * {{- . -}} Connection) CloseAndLogErrors() {
	c.CloseStream()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil && c.log != nil {
			c.log.Printf("gRPC Close() failure: %s", err)
		}
	}
	c.conn = nil
}

// Connect connects, a possibly half-connected {{ . -}} Connection,
// and return a new connection.
func (c * {{- . -}} Connection) Connect(address string, declaration * {{- . -}} Declaration, opts ...grpc.DialOption) (res * {{- . -}} Connection, err error) {
	defer func() {
		if err == nil {
			return
		}
		res.CloseAndLogErrors()
	}()
	res = & {{- . -}} Connection{
		log: c.log,
	}
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

	res.stream, err = client. {{- . -}} (context.Background(), DefaultCallOptions...)
	if err != nil {
		return
	}
	res.acknowledge, err =  res.Send(& {{- . -}} UpStream{
		Declaration: declaration,
	})
	return
}

// ConnectAsync Connects a posibly half-connected {{ . -}} Connection
// asynchronously.
func (c * {{- . -}} Connection) ConnectAsync(address string,
	declaration * {{- . -}} Declaration,
	opts ...grpc.DialOption) (<-chan * {{- . -}} Connection, <-chan error) {

	errors := make(chan error)
	connections := make(chan * {{- . -}} Connection)
	declaration = deepcopy.MustAnything(declaration).(* {{- . -}} Declaration)

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

{{ end }}

`))

func execute() error {
	output_path := "olympuspb/client_utils.go"
	fmt.Printf("Writing '%s'\n", output_path)
	f, err := os.Create(output_path)
	if err != nil {
		return err
	}
	defer f.Close()
	return client_util_templates.Execute(f, clientUtilTemplateData{
		Clients: []string{"Zone", "Tracking"},
	})
}
