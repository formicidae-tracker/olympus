package api

import (
	"context"
	"io"
	"net"
	"testing"

	"google.golang.org/grpc"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type olympusStub struct {
	UnimplementedOlympusServer
}

func (o *olympusStub) Tracking(stream Olympus_TrackingServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = stream.Send(&TrackingDownStream{})
		if err != nil {
			return err
		}
	}
}

type ClientConnectionSuite struct {
	server  *grpc.Server
	stub    *olympusStub
	errors  chan error
	address string
}

var _ = Suite(&ClientConnectionSuite{})

func (s *ClientConnectionSuite) SetUpTest(c *C) {
	s.address = "localhost:12345"
	s.stub = &olympusStub{}
	s.errors = make(chan error)
	s.server = grpc.NewServer()
	RegisterOlympusServer(s.server, s.stub)
	l, err := net.Listen("tcp", s.address)
	c.Assert(err, IsNil)
	go func() {
		s.errors <- s.server.Serve(l)
		close(s.errors)
	}()
}

func (s *ClientConnectionSuite) TearDownTest(c *C) {
	s.server.GracefulStop()
	err, ok := <-s.errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)
	err, ok = <-s.errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)
}

func (s *ClientConnectionSuite) TestConnect(c *C) {
	ch := ConnectTracking(context.Background(), s.address, &TrackingDeclaration{})
	res, ok := <-ch
	c.Assert(res, Not(IsNil))
	c.Assert(ok, Equals, true)
	_, ok = <-ch
	c.Check(ok, Equals, false)

	c.Assert(res.Connection, Not(IsNil))
	c.Check(res.Error, IsNil)
	c.Check(res.Connection.Close(), IsNil)
}
