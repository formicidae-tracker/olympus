package main

import (
	"io"

	"github.com/formicidae-tracker/olympus/olympuspb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OlympusGRPCWrapper Olympus

type serverStream[UpStream any, DownStream any] interface {
	Recv() (*UpStream, error)
	Send(*DownStream) error
}

func readAll[UpStream any, DownStream any](
	stream serverStream[UpStream, DownStream],
	messages chan<- *UpStream,
	errors chan<- error) {

	defer close(messages)
	defer close(errors)
	for {
		m, err := stream.Recv()
		if err != nil {
			errors <- err
		} else {
			messages <- m
		}
	}
}

func serveLoop[UpStream any, DownStream any](
	stream serverStream[UpStream, DownStream],
	handleMessage func(*UpStream) (*DownStream, error),
	finish *<-chan struct{}) error {

	messages := make(chan *UpStream)
	errors := make(chan error)

	go readAll(stream, messages, errors)

	for {
		select {
		case <-(*finish):
			return nil
		case err := <-errors:
			if err == io.EOF {
				return nil
			}
			return err
		case m := <-messages:
			out, err := handleMessage(m)
			if err != nil {
				return err
			}
			if out == nil {
				break
			}
			err = stream.Send(out)
			if err != nil {
				return err
			}
		}
	}
}

func mapError(err error) error {
	switch err.(type) {
	case AlreadyExistError:
		return status.Error(codes.AlreadyExists, err.Error())
	case ZoneNotFoundError:
		return status.Error(codes.NotFound, err.Error())
	case HostNotFoundError:
		return status.Error(codes.NotFound, err.Error())
	case ClosedOlympusServerError:
		return status.Error(codes.Internal, err.Error())
	default:
		return err
	}
}

func (o *OlympusGRPCWrapper) Zone(stream olympuspb.Olympus_ZoneServer) (err error) {
	var z ZoneLogger = nil
	var finish <-chan struct{} = nil

	defer func() {
		if z == nil {
			return
		}
		graceful := err != nil
		(*Olympus)(o).UnregisterZone(z.ZoneIdentifier(), graceful)
	}()
	ack := &olympuspb.ZoneDownStream{}
	handleMessage := func(m *olympuspb.ZoneUpStream) (*olympuspb.ZoneDownStream, error) {

		if z == nil {
			if m.Declaration == nil {
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain ZoneDeclaration")
			}
			var err error
			z, finish, err = (*Olympus)(o).RegisterZone(m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}
		}

		if m.Target != nil {
			z.PushTarget(m.Target)
		}
		if len(m.Alarms) > 0 {
			z.PushAlarms(m.Alarms)
		}
		if len(m.Reports) > 0 {
			z.PushReports(m.Reports)
		}
		return ack, nil
	}

	return serveLoop[olympuspb.ZoneUpStream, olympuspb.ZoneDownStream](stream, handleMessage, &finish)
}

func (o *OlympusGRPCWrapper) Tracking(stream olympuspb.Olympus_TrackingServer) (err error) {
	var t TrackingLogger = nil
	var finish <-chan struct{} = nil
	hostname := ""
	defer func() {
		if t == nil {
			return
		}
		graceful := err != nil
		(*Olympus)(o).UnregisterTracker(hostname, graceful)
	}()
	ack := &olympuspb.TrackingDownStream{}
	handleMessage := func(m *olympuspb.TrackingUpStream) (*olympuspb.TrackingDownStream, error) {
		if t == nil {
			if m.Declaration == nil {
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain TrackingDeclaration")
			}
			var err error
			t, finish, err = (*Olympus)(o).RegisterTracker(m.Declaration)
			hostname = m.Declaration.Hostname
			if err != nil {
				return nil, mapError(err)
			}
		}
		return ack, nil
	}

	return serveLoop[olympuspb.TrackingUpStream, olympuspb.TrackingDownStream](stream, handleMessage, &finish)
}
