package olympus

import (
	"context"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OlympusGRPCWrapper Olympus

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

func (o *OlympusGRPCWrapper) SubscriptionContext() context.Context {
	return (*Olympus)(o).subscriptionContext
}

func (o *OlympusGRPCWrapper) Log(ctx context.Context) *logrus.Entry {
	return (*Olympus)(o).log.WithField("domain", "gRPC").WithContext(ctx)
}

func (o *OlympusGRPCWrapper) Climate(stream api.Olympus_ClimateServer) (err error) {
	var subscription *GrpcSubscription[ClimateLogger] = nil
	ctx := api.WithTelemetry(o.SubscriptionContext(), "fort.olympus.Olympus/Climate")

	defer func() {
		if subscription == nil {
			return
		}

		if err != nil {
			o.Log(ctx).WithError(err).Error("climate stream error")
		}

		graceful := err == nil
		(*Olympus)(o).UnregisterClimate(ctx, subscription.object.Host(), subscription.object.ZoneName(), graceful)
	}()

	ack := &api.ClimateDownStream{}
	handler := func(mCtx context.Context, m *api.ClimateUpStream) (*api.ClimateDownStream, error) {
		var confirmation *api.ClimateDownStream
		if subscription == nil {
			ctx = mCtx

			if m.Declaration == nil {
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain ZoneDeclaration")
			}
			var err error

			subscription, err = (*Olympus)(o).RegisterClimate(ctx, m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}

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
			subscription.alarmLogger.PushAlarms(m.Alarms, "climate")
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

	return api.ServerLoop[*api.ClimateUpStream, *api.ClimateDownStream](
		ctx, stream, handler)
}

func (o *OlympusGRPCWrapper) Tracking(stream api.Olympus_TrackingServer) (err error) {
	var subscription *GrpcSubscription[TrackingLogger] = nil
	ctx := api.WithTelemetry(o.SubscriptionContext(), "fort.olympus.Olympus/Tracking")

	hostname := ""
	defer func() {
		if subscription == nil {
			return
		}

		if err != nil {
			o.Log(ctx).WithError(err).Error("tracking stream error")
		}

		graceful := err == nil
		(*Olympus)(o).UnregisterTracker(ctx, hostname, graceful)
	}()
	ack := &api.TrackingDownStream{}

	handler := func(mCtx context.Context, m *api.TrackingUpStream) (*api.TrackingDownStream, error) {
		if subscription == nil {
			ctx = mCtx
			if m.Declaration == nil {
				return nil, status.Error(codes.InvalidArgument, "first message of stream must contain TrackingDeclaration")
			}
			var err error
			subscription, err = (*Olympus)(o).RegisterTracking(ctx, m.Declaration)
			if err != nil {
				return nil, mapError(err)
			}
			hostname = m.Declaration.Hostname
		}

		if len(m.Alarms) > 0 {
			subscription.alarmLogger.PushAlarms(m.Alarms, "tracking")
			subscription.NotifyAlarms(m.Alarms)
		}

		if m.DiskStatus != nil {
			subscription.object.PushDiskStatus(m.DiskStatus)
		}

		return ack, nil
	}

	return api.ServerLoop[*api.TrackingUpStream, *api.TrackingDownStream](
		ctx, stream, handler)
}

func (o *OlympusGRPCWrapper) SendAlarm(ctx context.Context, update *api.AlarmUpdate) (*empty.Empty, error) {
	go (*Olympus)(o).NotifyAlarm(ctx, "fake.zone", update)
	return &empty.Empty{}, nil
}
