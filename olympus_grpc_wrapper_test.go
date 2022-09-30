package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/proto"
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

func (s *GRPCSuite) TestNothingHappens(c *C) {
	conn, err := grpc.Dial("localhost:12345", proto.DefaultDialOptions...)
	c.Assert(err, IsNil)
	defer conn.Close()
}

func (s *GRPCSuite) TestEndToEnd(c *C) {
	conn, err := grpc.Dial("localhost:12345", proto.DefaultDialOptions...)
	c.Assert(err, IsNil)
	defer conn.Close()

	client := proto.NewOlympusClient(conn)

	stream, err := client.Zone(context.Background(), proto.DefaultCallOptions...)
	c.Assert(err, IsNil)

	reports := []*proto.ClimateReport{
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

	target := &proto.ClimateTarget{
		Current: &proto.ClimateState{
			Name:         "box",
			Temperature:  newInitialized[float32](23.0),
			Humidity:     newInitialized[float32](55.0),
			Wind:         newInitialized[float32](100.0),
			VisibleLight: newInitialized[float32](0.0),
			UvLight:      newInitialized[float32](0.0),
		},
	}

	lastReports := reports[len(reports)-1]

	c.Check(stream.Send(&proto.ZoneUpStream{
		Declaration: &proto.ZoneDeclaration{
			Host: "somehost",
			Name: "box",
		},
		Reports: reports[4:],
		Target:  target,
	}), IsNil)

	c.Check(stream.Send(&proto.ZoneUpStream{
		Reports: reports[:4],
	}), IsNil)

	for s.o.ZoneIsRegistered("somehost", "box") == false {
		time.Sleep(1 * time.Millisecond)
	}

	for {
		series, err := s.o.GetClimateTimeSerie("somehost", "box", "")
		c.Assert(err, IsNil)
		fmt.Println(len(series.Humidity))
		if len(series.Humidity) == 5 {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	report, err := s.o.GetZoneReport("somehost", "box")
	if c.Check(err, IsNil) == true {
		c.Check(report, DeepEquals, ZoneReport{
			Host: "somehost",
			Name: "box",
			Climate: &ZoneClimateReport{
				Temperature: &lastReports.Temperatures[0],
				Humidity:    lastReports.Humidity,
				Current:     deepcopy.MustAnything(target.Current).(*proto.ClimateState),
			},
			Alarms: []AlarmReport{},
		})
	}

	c.Check(stream.CloseSend(), IsNil)

	for s.o.ZoneIsRegistered("somehost", "box") == true {
		time.Sleep(1 * time.Millisecond)
	}

}
