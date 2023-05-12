package main

import (
	"io"

	"github.com/formicidae-tracker/olympus/api"
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
	case NoClimateRunningError:
		return status.Error(codes.NotFound, err.Error())
	case NoTrackingRunningError:
		return status.Error(codes.NotFound, err.Error())
	case ZoneNotFoundError:
		return status.Error(codes.NotFound, err.Error())
	case ClosedOlympusServerError:
		return status.Error(codes.Internal, err.Error())
	default:
		return err
	}
}

func (o *OlympusGRPCWrapper) Climate(stream api.Olympus_ClimateServer) (err error) {
	var subscription *GrpcSubscription[ClimateLogger] = nil
	var finish <-chan struct{} = nil

	defer func() {
		if subscription == nil {
			return
		}
		graceful := err != nil
		(*Olympus)(o).UnregisterClimate(subscription.object.Host(), subscription.object.ZoneName(), graceful)
	}()
	ack := &api.ClimateDownStream{}
	handleMessage := func(m *api.ClimateUpStream) (*api.ClimateDownStream, error) {
		var confirmation *api.ClimateDownStream
		if subscription == nil {
			if m.Declaration == nil {
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain ZoneDeclaration")
			}
			var err error
			subscription, err = (*Olympus)(o).RegisterClimate(m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}
			finish = subscription.finish

			confirmation = &api.ClimateDownStream{
				RegistrationConfirmation: &api.ClimateRegistrationConfirmation{
					PageSize: int32(BackLogPageSize),
				},
			}
		}

		if m.Target != nil {
			subscription.object.PushTarget(m.Target)
		}
		if len(m.Alarms) > 0 {
			subscription.alarmLogger.PushAlarms(m.Alarms)
		}
		if len(m.Reports) > 0 {
			subscription.object.PushReports(m.Reports)
		}

		if confirmation != nil {
			return confirmation, nil
		}

		return ack, nil
	}

	return serveLoop[api.ClimateUpStream, api.ClimateDownStream](stream, handleMessage, &finish)
}

func (o *OlympusGRPCWrapper) Tracking(stream api.Olympus_TrackingServer) (err error) {
	var subscription *GrpcSubscription[TrackingLogger] = nil
	var finish <-chan struct{} = nil
	hostname := ""
	defer func() {
		if subscription == nil {
			return
		}
		graceful := err != nil
		(*Olympus)(o).UnregisterTracker(hostname, graceful)
	}()
	ack := &api.TrackingDownStream{}
	handleMessage := func(m *api.TrackingUpStream) (*api.TrackingDownStream, error) {
		if subscription == nil {
			if m.Declaration == nil {
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain TrackingDeclaration")
			}
			var err error
			subscription, err = (*Olympus)(o).RegisterTracking(m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}
			hostname = m.Declaration.Hostname
			finish = subscription.finish

		}

		if len(m.Alarms) > 0 {
			subscription.alarmLogger.PushAlarms(m.Alarms)
		}

		return ack, nil
	}

	return serveLoop[api.TrackingUpStream, api.TrackingDownStream](stream, handleMessage, &finish)
}
