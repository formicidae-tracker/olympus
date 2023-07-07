package api

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	. "gopkg.in/check.v1"
)

type ClientTaskSuite struct {
	ctrl    *gomock.Controller
	olympus *MockOlympusServer

	server *grpc.Server
	errors chan error
}

var _ = Suite(&ClientTaskSuite{})

func (s *ClientTaskSuite) SetUpTest(c *C) {
	s.ctrl = gomock.NewController(NeverFatalReporter{c})
	s.olympus = NewMockOlympusServer(s.ctrl)

	s.server = grpc.NewServer(DefaultServerOptions...)

	s.errors = make(chan error)
	RegisterOlympusServer(s.server, s.olympus)
	l, err := net.Listen("tcp", testAddress)
	if c.Check(err, IsNil) == false {
		close(s.errors)
	} else {
		go func() {
			s.errors <- s.server.Serve(l)
			close(s.errors)
		}()
	}
}

func (s *ClientTaskSuite) TearDownTest(c *C) {

	s.server.GracefulStop()

	err, ok := <-s.errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)

	err, ok = <-s.errors
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)

	s.ctrl.Finish()
}

func getOrTimeout[T any](ch <-chan T, timeout time.Duration, c *C) (T, bool) {

	select {
	case res, ok := <-ch:
		return res, ok
	case <-time.After(timeout):
		c.Errorf("could not fetch from channel after %s", timeout)
		var res T
		return res, false
	}

}

func (s *ClientTaskSuite) TestEndsWithAnEOF(c *C) {

	task := NewClimateTask(context.Background(),
		testAddress, &ClimateDeclaration{})
	done := make(chan struct{})
	defer func() {
		<-done
	}()
	s.olympus.EXPECT().
		Climate(gomock.Any()).
		DoAndReturn(func(s Olympus_ClimateServer) error {
			defer close(done)
			err := acknowledgeAll[ClimateUpStream, ClimateDownStream](s,
				&ClimateDownStream{
					RegistrationConfirmation: &ClimateRegistrationConfirmation{
						SendBacklogs: true,
					},
				}, &ClimateDownStream{})
			c.Check(err, Equals, io.EOF)
			return err
		})

	errors := make(chan error)
	go func() {
		defer close(errors)
		errors <- task.Run()
	}()

	timeout := 50 * time.Millisecond
	confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(confirmation.Confirmation, Not(IsNil))
	c.Check(confirmation.Error, IsNil)
	c.Check(ok, Equals, true)

	req := task.Request(&ClimateUpStream{})
	res, ok := getOrTimeout(req, timeout, c)

	c.Check(res.Message, Not(IsNil))
	c.Check(res.Error, IsNil)
	c.Check(ok, Equals, true)

	_, ok = getOrTimeout(req, timeout, c)
	c.Check(ok, Equals, false)
	task.stop()

	_, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, false)

	err, ok := getOrTimeout(errors, timeout, c)
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)

	_, ok = getOrTimeout(errors, timeout, c)
	c.Check(ok, Equals, false)

}

func (s *ClientTaskSuite) TestReconnectionOnError(c *C) {
	task := NewClimateTask(context.Background(),
		testAddress, &ClimateDeclaration{})
	done := make(chan struct{})
	defer func() { <-done }()

	gomock.InOrder(
		s.olympus.EXPECT().
			Climate(gomock.Any()).
			DoAndReturn(func(s Olympus_ClimateServer) error {
				_, err := s.Recv()
				if err != nil {
					return err
				}
				return grpc.Errorf(codes.InvalidArgument, "you are a nasty buoy")
			}),
		s.olympus.EXPECT().
			Climate(gomock.Any()).
			DoAndReturn(func(s Olympus_ClimateServer) error {
				defer close(done)
				err := acknowledgeAll[ClimateUpStream, ClimateDownStream](s,
					&ClimateDownStream{
						RegistrationConfirmation: &ClimateRegistrationConfirmation{
							SendBacklogs: true,
						},
					}, &ClimateDownStream{})
				c.Check(err, Equals, io.EOF)
				return err
			}),
	)

	errors := make(chan error)
	go func() {
		defer close(errors)
		errors <- task.Run()
	}()

	timeout := 50 * time.Millisecond
	confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(confirmation.Confirmation, IsNil)
	c.Check(confirmation.Error, ErrorMatches, ".*InvalidArgument.*you are a nasty buoy")
	c.Check(ok, Equals, true)

	confirmation, ok = getOrTimeout(task.Confirmations(), timeout, c)
	if c.Check(confirmation.Confirmation, Not(IsNil)) == true {
		c.Check(confirmation.Confirmation.RegistrationConfirmation, Not(IsNil))
	}
	c.Check(confirmation.Error, IsNil)
	c.Check(ok, Equals, true)

	task.stop()

	confirmation, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(confirmation.Confirmation, IsNil)
	c.Check(confirmation.Error, IsNil)
	c.Check(ok, Equals, false)

}

func errorsOnFirstRequest[Up, Down any](s server[Up, Down], conf *Down) error {
	_, err := s.Recv()
	if err != nil {
		return err
	}
	err = s.Send(conf)
	if err != nil {
		return err
	}
	_, err = s.Recv()
	if err != nil {
		return err
	}
	return grpc.Errorf(codes.InvalidArgument, "you are a not so nasty buoy")
}

func (s *ClientTaskSuite) TestCloseAndReconnectOnError(c *C) {
	task := NewClimateTask(context.Background(),
		testAddress, &ClimateDeclaration{})
	done := make(chan struct{})
	defer func() { <-done }()

	conf := &ClimateDownStream{
		RegistrationConfirmation: &ClimateRegistrationConfirmation{
			SendBacklogs: true,
		},
	}

	ack := &ClimateDownStream{}

	gomock.InOrder(
		s.olympus.EXPECT().
			Climate(gomock.Any()).
			DoAndReturn(func(s Olympus_ClimateServer) error {
				return errorsOnFirstRequest[ClimateUpStream, ClimateDownStream](s,
					conf)
			}),
		s.olympus.EXPECT().
			Climate(gomock.Any()).
			DoAndReturn(func(s Olympus_ClimateServer) error {
				defer close(done)
				err := acknowledgeAll[ClimateUpStream, ClimateDownStream](s,
					conf, ack)
				c.Check(err, Equals, io.EOF)
				return err
			}),
	)

	errors := make(chan error)
	go func() {
		defer close(errors)
		errors <- task.Run()
	}()

	timeout := 50 * time.Millisecond
	confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
	if c.Check(confirmation.Confirmation, Not(IsNil)) == true {
		c.Check(confirmation.Confirmation.RegistrationConfirmation, Not(IsNil))
	}
	c.Check(confirmation.Error, IsNil)
	c.Check(ok, Equals, true)

	resp := task.Request(&ClimateUpStream{})
	res, ok := getOrTimeout(resp, timeout, c)
	c.Check(res.Message, IsNil)
	c.Check(res.Error, ErrorMatches, ".*InvalidArgument.*you are a not so nasty buoy")
	c.Check(ok, Equals, true)

	confirmation, ok = getOrTimeout(task.Confirmations(), timeout, c)
	if c.Check(confirmation.Confirmation, Not(IsNil)) == true {
		c.Check(confirmation.Confirmation.RegistrationConfirmation, Not(IsNil))
	}
	c.Check(confirmation.Error, IsNil)
	c.Check(ok, Equals, true)

	task.stop()

	confirmation, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(confirmation.Confirmation, IsNil)
	c.Check(confirmation.Error, IsNil)
	c.Check(ok, Equals, false)

}
