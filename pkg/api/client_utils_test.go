package api

import (
	"io"
	"net"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
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

type StarConnectionSuite struct {
	server  *grpc.Server
	stub    *olympusStub
	errors  chan error
	address string
}

var _ = Suite(&StarConnectionSuite{})

func (s *StarConnectionSuite) SetUpTest(c *C) {
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

func (s *StarConnectionSuite) TearDownTest(c *C) {
	s.server.GracefulStop()
	err, ok := <-s.errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)
	err, ok = <-s.errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)
}

func (s *StarConnectionSuite) TestConnect(c *C) {
	conn, err := ConnectTracking(nil, s.address, &TrackingDeclaration{}, nil)
	c.Assert(conn, Not(IsNil))
	c.Check(err, IsNil)
	conn.CloseAll(nil)
}

func (s *StarConnectionSuite) TestConnectAsync(c *C) {
	logger, hook := test.NewNullLogger()
	defer hook.Reset()
	connections, errors := ConnectTrackingAsync(nil, s.address, &TrackingDeclaration{}, logger.WithField("group", "gRPC"))

	conn, ok := <-connections
	c.Assert(conn, Not(IsNil))
	c.Assert(ok, Equals, true)

	_, ok = <-connections
	c.Check(ok, Equals, false)

	err, ok := <-errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)

	conn.CloseAll(nil)
}
