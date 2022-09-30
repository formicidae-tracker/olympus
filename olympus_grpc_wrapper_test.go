package main

import (
	"net"

	"github.com/formicidae-tracker/olympus/proto"
	"google.golang.org/grpc"
	. "gopkg.in/check.v1"
)

type GRPCSuite struct {
	o        *Olympus
	server   *grpc.Server
	shutdown chan struct{}
	done     chan error
}

func (s *GRPCSuite) initialize() error {
	var err error
	s.o, err = NewOlympus("")
	if err != nil {
		return err
	}

	s.server = grpc.NewServer(proto.DefaultServerOptions...)
	proto.RegisterOlympusServer(s.server, (*OlympusGRPCWrapper)(s.o))

	s.shutdown = make(chan struct{})
	s.done = make(chan error)
	return nil
}

func (s *GRPCSuite) serveAndListen() error {
	lis, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		return err
	}

	go func() {
		s.done <- s.server.Serve(lis)
		close(s.done)
	}()
	go func() {
		<-s.shutdown
		s.server.GracefulStop()
	}()
	return nil
}

var _ = Suite(&GRPCSuite{})

func (s *GRPCSuite) SetUpTest(c *C) {
	c.Assert(s.initialize(), IsNil)
	c.Assert(s.serveAndListen(), IsNil)
}

func (s *GRPCSuite) TearDownTest(c *C) {
	close(s.shutdown)
	c.Check(s.o.Close(), IsNil)
	err, ok := <-s.done
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)
	err, ok = <-s.done
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)
}

func (s *GRPCSuite) TestNothinHappens(c *C) {
}
