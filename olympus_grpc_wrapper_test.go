package main

import (
	"context"
	"net"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/olympuspb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	s.server = grpc.NewServer(olympuspb.DefaultServerOptions...)
	olympuspb.RegisterOlympusServer(s.server, (*OlympusGRPCWrapper)(s.o))

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

func (s *GRPCSuite) TestNothingHappens(c *C) {
	conn, err := grpc.Dial("localhost:12345", olympuspb.DefaultDialOptions...)
	c.Assert(err, IsNil)
	defer conn.Close()
}

func connectZone(c *C) (olympuspb.Olympus_ZoneClient, func(), error) {
	conn, err := grpc.Dial("localhost:12345", olympuspb.DefaultDialOptions...)
	if err != nil {
		return nil, func() {}, err
	}

	client := olympuspb.NewOlympusClient(conn)
	stream, err := client.Zone(context.Background(), olympuspb.DefaultCallOptions...)
	if err != nil {
		return nil, func() { c.Check(conn.Close(), IsNil) }, err
	}
	return stream, func() {
		c.Check(stream.CloseSend(), IsNil)
		c.Check(conn.Close(), IsNil)
	}, nil

}

func connectTracking(c *C) (olympuspb.Olympus_TrackingClient, func(), error) {
	conn, err := grpc.Dial("localhost:12345", olympuspb.DefaultDialOptions...)
	if err != nil {
		return nil, func() {}, err
	}

	client := olympuspb.NewOlympusClient(conn)
	stream, err := client.Tracking(context.Background(), olympuspb.DefaultCallOptions...)
	if err != nil {
		return nil, func() { c.Check(conn.Close(), IsNil) }, err
	}
	return stream, func() {
		c.Check(stream.CloseSend(), IsNil)
		c.Check(conn.Close(), IsNil)
	}, nil
}

func (s *GRPCSuite) TestEndToEnd(c *C) {
	stream, cleanUp, err := connectZone(c)
	defer cleanUp()
	c.Assert(err, IsNil)

	reports := []*olympuspb.ClimateReport{
		{
			Time:         timestamppb.New(time.Time{}.Add(500 * time.Millisecond)),
			Humidity:     newInitialized[float32](55.3),
			Temperatures: []float32{22.1, 23.04, 22.97, 21.97},
		},
		{
			Time:         timestamppb.New(time.Time{}.Add(1000 * time.Millisecond)),
			Humidity:     newInitialized[float32](55.3),
			Temperatures: []float32{22.1, 23.04, 22.97, 21.97},
		},
		{
			Time:         timestamppb.New(time.Time{}.Add(1500 * time.Millisecond)),
			Humidity:     newInitialized[float32](55.3),
			Temperatures: []float32{22.1, 23.04, 22.97, 21.97},
		},
		{
			Time:         timestamppb.New(time.Time{}.Add(2000 * time.Millisecond)),
			Humidity:     newInitialized[float32](55.3),
			Temperatures: []float32{22.1, 23.04, 22.97, 21.97},
		},
		{
			Time:         timestamppb.New(time.Time{}.Add(2500 * time.Millisecond)),
			Humidity:     newInitialized[float32](55.2),
			Temperatures: []float32{22.8, 23.1, 22.9, 21.9},
		},
	}

	target := &olympuspb.ClimateTarget{
		Current: &olympuspb.ClimateState{
			Name:         "box",
			Temperature:  newInitialized[float32](23.0),
			Humidity:     newInitialized[float32](55.0),
			Wind:         newInitialized[float32](100.0),
			VisibleLight: newInitialized[float32](0.0),
			UvLight:      newInitialized[float32](0.0),
		},
	}

	lastReports := reports[len(reports)-1]

	c.Check(stream.Send(&olympuspb.ZoneUpStream{
		Declaration: &olympuspb.ZoneDeclaration{
			Host: "somehost",
			Name: "box",
		},
		Reports: reports[4:],
		Target:  target,
	}), IsNil)
	_, err = stream.Recv()
	c.Check(err, IsNil)

	c.Check(stream.Send(&olympuspb.ZoneUpStream{
		Reports: reports[:4],
	}), IsNil)
	_, err = stream.Recv()
	c.Check(err, IsNil)

	report, err := s.o.GetZoneReport("somehost", "box")
	if c.Check(err, IsNil) == true {
		c.Check(report, DeepEquals, ZoneReport{
			Host: "somehost",
			Name: "box",
			Climate: &ZoneClimateReport{
				Temperature: &lastReports.Temperatures[0],
				Humidity:    lastReports.Humidity,
				Current:     deepcopy.MustAnything(target.Current).(*olympuspb.ClimateState),
			},
			Alarms: []AlarmReport{},
		})
	}
	c.Check(stream.CloseSend(), IsNil)

	for s.o.ZoneIsRegistered("somehost", "box") == true {
		time.Sleep(1 * time.Millisecond)
	}

}

func (s *GRPCSuite) TestDoubleZoneRegistrationError(c *C) {
	streams := []olympuspb.Olympus_ZoneClient{nil, nil}

	for i := range streams {
		stream, cleanUp, err := connectZone(c)
		defer cleanUp()
		c.Assert(err, IsNil)
		streams[i] = stream
	}
	declaration := &olympuspb.ZoneUpStream{
		Declaration: &olympuspb.ZoneDeclaration{Host: "somehost", Name: "box"},
	}
	c.Check(streams[0].Send(declaration), IsNil)
	_, err := streams[0].Recv()
	c.Check(err, IsNil)

	c.Check(streams[1].Send(declaration), IsNil)
	m, err := streams[1].Recv()
	c.Check(m, IsNil)
	c.Check(err, ErrorMatches, `rpc error: code = AlreadyExists desc = zone 'somehost.box' is already registered`)
}

func (s *GRPCSuite) TestLackOfZonRegistrationError(c *C) {
	stream, cleanUp, err := connectZone(c)
	defer cleanUp()
	c.Assert(err, IsNil)

	c.Check(stream.Send(&olympuspb.ZoneUpStream{}), IsNil)
	m, err := stream.Recv()
	c.Check(m, IsNil)
	c.Check(err, ErrorMatches, `rpc error: code = InvalidArgument desc = first message of stream must contain ZoneDeclaration`)
}
