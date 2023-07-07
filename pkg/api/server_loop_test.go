package api

import (
	"context"
	"net"
	sync "sync"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	. "gopkg.in/check.v1"
)

type ServerLoopSuite struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	ctrl    *gomock.Controller
	olympus *MockOlympusServer

	server                         *grpc.Server
	serverErrors, connectionErrors chan error
}

var _ = Suite(&ServerLoopSuite{})

type handler struct {
	declarationReceived bool
}

func (h *handler) handle(ctx context.Context, m *TrackingUpStream) (*TrackingDownStream, error) {
	if h.declarationReceived == false && m.Declaration == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing declaration")
	} else if len(m.Alarms) > 0 {
		return nil, grpc.Errorf(codes.FailedPrecondition, "not yet implemented")
	}
	h.declarationReceived = true

	return &TrackingDownStream{}, nil
}

func (s *ServerLoopSuite) SetUpTest(c *C) {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.ctrl = gomock.NewController(NeverFatalReporter{c})
	s.olympus = NewMockOlympusServer(s.ctrl)

	s.connectionErrors = make(chan error, 100)
	s.olympus.EXPECT().
		Tracking(gomock.Any()).
		DoAndReturn(func(se Olympus_TrackingServer) error {
			s.wg.Add(1)
			defer s.wg.Done()
			handler := &handler{}
			err := ServerLoop[*TrackingUpStream, *TrackingDownStream](s.ctx,
				se, handler.handle)
			s.connectionErrors <- err
			return err
		}).
		AnyTimes()

	s.server = grpc.NewServer(DefaultServerOptions...)
	RegisterOlympusServer(s.server, s.olympus)

	s.serverErrors = make(chan error)
	l, err := net.Listen("tcp", testAddress)
	if c.Check(err, IsNil) == false {
		close(s.serverErrors)
		return
	}

	go func() {
		defer close(s.serverErrors)
		s.serverErrors <- s.server.Serve(l)
	}()
}

func waitOrTimeout(wg *sync.WaitGroup, timeout time.Duration, c *C) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		c.Errorf("WaitGroup is not Done() after %s", timeout)
	}
}

func (s *ServerLoopSuite) TearDownTest(c *C) {
	s.server.GracefulStop()
	waitOrTimeout(&s.wg, 100*time.Millisecond, c)

	err, ok := getOrTimeout(s.serverErrors, 10*time.Millisecond, c)
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)
	_, ok = getOrTimeout(s.serverErrors, 10*time.Millisecond, c)
	c.Check(ok, Equals, false)

	s.ctrl.Finish()
}

func (s *ServerLoopSuite) TestDoNothingConnection(c *C) {
	task := NewTrackingTask(context.Background(), testAddress, &TrackingDeclaration{})

	errors := make(chan error)
	go func() {
		defer close(errors)
		errors <- task.Run()
	}()
	timeout := 50 * time.Millisecond

	confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, true)
	c.Check(confirmation.Confirmation, Not(IsNil))
	c.Check(confirmation.Error, IsNil)

	task.stop()

	_, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, false)

	err, ok := getOrTimeout(errors, timeout, c)
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)

	err, ok = getOrTimeout(s.connectionErrors, 10*time.Millisecond, c)
	c.Check(ok, Equals, true)
	c.Check(err, IsNil)
}

func (s *ServerLoopSuite) TestServerCancelable(c *C) {
	clientsContext, cancelClients := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer func() {
		cancelClients()
		waitOrTimeout(&wg, 100*time.Millisecond, c)
	}()

	tasks := make([]*ClientTask[*TrackingUpStream, *TrackingDownStream], 10)

	for i := range tasks {
		tasks[i] = NewTrackingTask(clientsContext,
			testAddress, &TrackingDeclaration{})
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Check(tasks[i].Run(), IsNil)
		}(i)
	}

	timeout := 100 * time.Millisecond
	for i, task := range tasks {
		comment := Commentf("task %d", i)
		confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
		c.Check(ok, Equals, true, comment)
		c.Check(confirmation.Confirmation, Not(IsNil), comment)
		c.Check(confirmation.Error, IsNil, comment)
	}

	s.cancel()
	waitOrTimeout(&s.wg, 100*time.Millisecond, c)

	for range tasks {
		err, ok := getOrTimeout(s.connectionErrors, 10*time.Millisecond, c)
		c.Check(ok, Equals, true)
		c.Check(err, IsNil)
	}
}

func (s *ServerLoopSuite) TestNoErrorOnCleanCloseStream(c *C) {
	task := NewTrackingTask(context.Background(), testAddress, &TrackingDeclaration{})

	done := make(chan struct{})
	defer func() { <-done }()
	go func() {
		defer close(done)
		c.Check(task.Run(), IsNil)
	}()

	timeout := 50 * time.Millisecond
	confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, true)
	c.Check(confirmation.Confirmation, Not(IsNil))
	c.Check(confirmation.Error, IsNil)

	for i := 0; i < 10; i++ {
		resp := task.Request(&TrackingUpStream{})
		res, ok := getOrTimeout(resp, timeout, c)
		c.Check(ok, Equals, true)
		c.Check(res.Error, IsNil)
		c.Check(res.Message, Not(IsNil))
	}

	task.stop()

	_, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, false)

	err, ok := getOrTimeout(s.connectionErrors, 10*time.Millisecond, c)
	c.Check(ok, Equals, true)
	c.Check(err, IsNil)
}

func (s *ServerLoopSuite) TestErrorsOnBadClientStream(c *C) {
	conn, err := grpc.Dial(testAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	c.Assert(err, IsNil)

	client := NewOlympusClient(conn)
	stream, err := client.Tracking(context.Background())
	c.Assert(err, IsNil)

	err = stream.Send(&TrackingUpStream{Declaration: &TrackingDeclaration{}})
	c.Assert(err, IsNil)
	_, err = stream.Recv()
	c.Assert(err, IsNil)

	c.Assert(conn.Close(), IsNil)
	err, ok := getOrTimeout(s.connectionErrors, 10*time.Millisecond, c)
	c.Check(ok, Equals, true)
	c.Check(err, ErrorMatches, ".*context canceled")

}

func (s *ServerLoopSuite) TestServerError(c *C) {
	task := NewTrackingTask(context.Background(), testAddress, &TrackingDeclaration{})

	done := make(chan struct{})
	defer func() { <-done }()
	go func() {
		defer close(done)
		c.Check(task.Run(), IsNil)
	}()

	timeout := 50 * time.Millisecond
	confirmation, ok := getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, true)
	c.Check(confirmation.Confirmation, Not(IsNil))
	c.Check(confirmation.Error, IsNil)

	resp := task.Request(&TrackingUpStream{
		Alarms: []*AlarmUpdate{{Identification: "foo"}},
	})

	res, ok := getOrTimeout(resp, timeout, c)
	c.Check(ok, Equals, true)
	c.Check(res.Message, IsNil)
	c.Check(res.Error, ErrorMatches, ".*FailedPrecondition.*not yet implemented")

	confirmation, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, true)
	c.Check(confirmation.Confirmation, Not(IsNil))
	c.Check(confirmation.Error, IsNil)

	task.stop()

	_, ok = getOrTimeout(task.Confirmations(), timeout, c)
	c.Check(ok, Equals, false)

	err, ok := getOrTimeout(s.connectionErrors, 10*time.Millisecond, c)
	c.Check(ok, Equals, true)
	c.Check(err, ErrorMatches, ".*FailedPrecondition.*not yet implemented")
}
