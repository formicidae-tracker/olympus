package olympus

import (
	"context"
	"io"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
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

	// The channel arenever closed as receiving io.EOF error should
	// simply close all goroutines.

	for {
		m, err := stream.Recv()
		if err != nil {
			errors <- err
			if err == io.EOF {
				// no more message to read after an EOF, listening
				// goroutine should stop too.
				return
			}
		} else {
			messages <- m
		}
	}
}

func serveLoop[UpStream any, DownStream any](
	stream serverStream[UpStream, DownStream],
	handleMessage func(*UpStream) (*DownStream, error),
	entry *logrus.Entry,
	ctx context.Context) error {

	messages := make(chan *UpStream)
	errors := make(chan error)

	go readAll(stream, messages, errors)

	for {
		select {
		case <-ctx.Done():
			// we were asked to stop the connection
			return nil
		case err := <-errors:
			if err == io.EOF {
				// we received an EOF : Simply end loop
				return nil
			}
			return err
		case m := <-messages:
			out, err := handleMessage(m)
			if err != nil {
				return err
			}
			if out == nil {
				// wait for a new message
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
	case UnexpectedStreamServerError:
		return status.Error(codes.InvalidArgument, err.Error())
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

func (o *OlympusGRPCWrapper) Context() context.Context {
	return (*Olympus)(o).subscriptionContext
}

func (o *OlympusGRPCWrapper) Climate(stream api.Olympus_ClimateServer) (err error) {
	var subscription *GrpcSubscription[ClimateLogger] = nil
	defer func() {
		if subscription == nil {
			return
		}
		graceful := err != nil && err != io.EOF && err != context.Canceled
		(*Olympus)(o).UnregisterClimate(subscription.object.Host(), subscription.object.ZoneName(), graceful)
	}()

	entry := (*Olympus)(o).log

	ack := &api.ClimateDownStream{}
	handleMessage := func(m *api.ClimateUpStream) (*api.ClimateDownStream, error) {
		var confirmation *api.ClimateDownStream
		if subscription == nil {
			if m.Declaration == nil {
				entry.Errorf("received climate stream without declaration")
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain ZoneDeclaration")
			}
			var err error
			entry = entry.WithField("zone", ZoneIdentifier(m.Declaration.Host, m.Declaration.Name))

			subscription, err = (*Olympus)(o).RegisterClimate(m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}

			confirmation = &api.ClimateDownStream{
				RegistrationConfirmation: &api.ClimateRegistrationConfirmation{
					PageSize: int32(BackLogPageSize),
				},
			}
		}

		entry.WithField("message", m).Tracef("received message")

		if m.Target != nil {
			subscription.object.PushTarget(m.Target)
		}
		if len(m.Alarms) > 0 {
			subscription.alarmLogger.PushAlarms(m.Alarms)
			if m.Backlog == false {
				subscription.NotifyAlarms(m.Alarms)
			}
		}
		if len(m.Reports) > 0 {
			subscription.object.PushReports(m.Reports)
		}

		if confirmation != nil {
			return confirmation, nil
		}

		return ack, nil
	}

	return serveLoop[api.ClimateUpStream, api.ClimateDownStream](stream, handleMessage, entry, o.Context())
}

func (o *OlympusGRPCWrapper) Tracking(stream api.Olympus_TrackingServer) (err error) {
	var subscription *GrpcSubscription[TrackingLogger] = nil
	hostname := ""
	defer func() {
		if subscription == nil {
			return
		}
		graceful := err != nil && err != io.EOF && err != context.Canceled
		(*Olympus)(o).UnregisterTracker(hostname, graceful)
	}()
	ack := &api.TrackingDownStream{}

	entry := (*Olympus)(o).log

	handleMessage := func(m *api.TrackingUpStream) (*api.TrackingDownStream, error) {
		if subscription == nil {
			if m.Declaration == nil {
				entry.Errorf("received tracking stream without a declaration")
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain TrackingDeclaration")
			}
			entry = entry.WithField("zone", m.Declaration.Hostname+".box")
			var err error
			subscription, err = (*Olympus)(o).RegisterTracking(m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}
			hostname = m.Declaration.Hostname
		}

		entry.WithField("message", m).Tracef("received message")

		if len(m.Alarms) > 0 {
			subscription.alarmLogger.PushAlarms(m.Alarms)
			subscription.NotifyAlarms(m.Alarms)
		}

		if m.DiskStatus != nil {
			subscription.object.PushDiskStatus(m.DiskStatus)
		}

		return ack, nil
	}

	return serveLoop[api.TrackingUpStream, api.TrackingDownStream](stream, handleMessage, entry, o.Context())
}

func (o *OlympusGRPCWrapper) SendAlarm(ctx context.Context, update *api.AlarmUpdate) (*empty.Empty, error) {
	go (*Olympus)(o).NotifyAlarm("fake.zone", update)
	return &empty.Empty{}, nil
}
