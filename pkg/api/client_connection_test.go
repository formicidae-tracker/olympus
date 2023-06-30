package api

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ClientConnectionSuite struct {
	server  *grpc.Server
	errors  chan error
	address string
	ctrl    *gomock.Controller
	olympus *MockOlympusServer
}

type server[Up, Down any] interface {
	Recv() (*Up, error)
	Send(*Down) error
}

func acknowledgeAll[Up, Down any](s server[Up, Down], conf *Down, ack *Down) error {
	toSend := conf
	for {
		_, err := s.Recv()
		if err != nil {
			return err
		}
		err = s.Send(toSend)
		toSend = ack
		if err != nil {
			return err
		}
	}
}

var _ = Suite(&ClientConnectionSuite{})

type NeverFatalReporter struct {
	c *C
}

func (r NeverFatalReporter) Errorf(format string, args ...any) {
	r.c.Errorf(format, args...)
}

func (r NeverFatalReporter) Fatalf(format string, args ...any) {
	r.c.Errorf(format, args...)
}

func (s *ClientConnectionSuite) SetUpTest(c *C) {

	s.ctrl = gomock.NewController(NeverFatalReporter{c})
	s.olympus = NewMockOlympusServer(s.ctrl)
	s.address = "localhost:12345"
	s.errors = make(chan error)
	s.server = grpc.NewServer()
	RegisterOlympusServer(s.server, s.olympus)
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
	s.ctrl.Finish()
}

func (s *ClientConnectionSuite) TestConnect(c *C) {
	done := make(chan struct{})
	defer func() { <-done }()
	s.olympus.EXPECT().
		Tracking(gomock.Any()).
		DoAndReturn(
			func(s Olympus_TrackingServer) error {
				err := acknowledgeAll[TrackingUpStream, TrackingDownStream](
					s,
					&TrackingDownStream{},
					&TrackingDownStream{})
				c.Check(err, ErrorMatches, "EOF")
				close(done)
				return err
			})

	ch := ConnectTracking(context.Background(), s.address, &TrackingDeclaration{})
	res, ok := <-ch
	c.Assert(res, Not(IsNil))
	c.Assert(ok, Equals, true)
	_, ok = <-ch
	c.Check(ok, Equals, false)

	c.Check(res.Error, IsNil)
	c.Assert(res.Connection, Not(IsNil))
	defer func() { c.Check(res.Connection.Close(), IsNil) }()
	c.Check(res.Connection.Established(), Equals, true)

	c.Check(res.Connection.Confirmation(), Not(IsNil))

	for i := 0; i < 5; i++ {
		ack, err := res.Connection.Send(&TrackingUpStream{})
		c.Check(ack, Not(IsNil))
		c.Check(err, IsNil)
	}

}

func (s *ClientConnectionSuite) TestComfirmation(c *C) {
	done := make(chan struct{})
	defer func() { <-done }()
	s.olympus.EXPECT().
		Climate(gomock.Any()).
		DoAndReturn(
			func(s Olympus_ClimateServer) error {
				err := acknowledgeAll[ClimateUpStream, ClimateDownStream](
					s,
					&ClimateDownStream{
						RegistrationConfirmation: &ClimateRegistrationConfirmation{
							SendBacklogs: true,
						},
					},
					&ClimateDownStream{})
				c.Check(err, ErrorMatches, "EOF")
				close(done)
				return err
			})

	ch := ConnectClimate(context.Background(), s.address, &ClimateDeclaration{})
	res, ok := <-ch
	c.Assert(res, Not(IsNil))
	c.Assert(ok, Equals, true)
	_, ok = <-ch
	c.Check(ok, Equals, false)

	c.Check(res.Error, IsNil)
	c.Assert(res.Connection, Not(IsNil))
	defer func() { c.Check(res.Connection.Close(), IsNil) }()
	c.Check(res.Connection.Established(), Equals, true)

	confirmation := res.Connection.Confirmation()
	c.Assert(confirmation, Not(IsNil))
	c.Assert(confirmation.RegistrationConfirmation, Not(IsNil))

	for i := 0; i < 5; i++ {
		ack, err := res.Connection.Send(&ClimateUpStream{})
		c.Check(ack, Not(IsNil))
		c.Check(err, IsNil)
	}

}
