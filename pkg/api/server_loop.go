package api

import (
	"context"
	"io"

	"github.com/formicidae-tracker/olympus/pkg/tm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type HandleFunc[Up, Down any] func(context.Context, Up) (Down, error)

type ServerStream[Up, Down metadated] interface {
	Recv() (Up, error)
	Send(Down) error
	Context() context.Context
}

type recvResult[Up any] struct {
	Message Up
	Error   error
}

func readAll[Up, Down metadated](ctx context.Context, ch chan<- recvResult[Up], s ServerStream[Up, Down]) {

	defer close(ch)

	for {
		m, err := s.Recv()

		ch <- recvResult[Up]{m, err}

		if err != nil {
			return
		}
	}
}

type serverTelemetry struct {
	serviceName string
	tracer      trace.Tracer
	propagator  propagation.TextMapPropagator
	links       []trace.Link
	linkedFirst bool
}

type telemetryKey string

var key telemetryKey = "telemetry"

func WithTelemetry(ctx context.Context, serviceName string) context.Context {
	if tm.Enabled() == false {
		return ctx
	}
	tlm := &serverTelemetry{
		serviceName: serviceName,
		tracer:      otel.GetTracerProvider().Tracer("github.com/formicidae-tracker/olympus/pkg/tm"),
		propagator:  otel.GetTextMapPropagator(),
	}
	return context.WithValue(ctx, key, tlm)
}

func getTelemetry(ctx context.Context) *serverTelemetry {
	res := ctx.Value(key)
	if res == nil {
		return nil
	}
	return res.(*serverTelemetry)
}

func (tm *serverTelemetry) linkWithGRPC(ctx context.Context) {
	grpcSpanContext := trace.SpanContextFromContext(ctx)
	if grpcSpanContext.IsValid() == false {
		return
	}
	tm.links = append(tm.links, trace.Link{SpanContext: grpcSpanContext})
}

func startSpan[Up metadated](tm *serverTelemetry, m Up) (context.Context, trace.Span) {
	if tm == nil {
		return context.Background(), nil
	}

	ctx := tm.propagator.Extract(context.Background(), textCarrier{m})

	ctx, span := tm.tracer.Start(ctx,
		tm.serviceName,
		trace.WithLinks(tm.links...),
	)

	if tm.linkedFirst == false {
		tm.links = append(tm.links, trace.Link{
			SpanContext: span.SpanContext(),
		})
		tm.linkedFirst = true
	}

	return ctx, span
}

func ServerLoop[Up, Down metadated](
	ctx context.Context,
	s ServerStream[Up, Down],
	handler HandleFunc[Up, Down]) error {

	recv := make(chan recvResult[Up])
	go func() {
		readAll(ctx, recv, s)
	}()

	tm := getTelemetry(ctx)
	if tm != nil {
		tm.linkWithGRPC(s.Context())
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case m, ok := <-recv:
			if ok == false {
				return nil
			}
			if m.Error != nil {
				if m.Error == io.EOF {
					return nil
				}
				return m.Error
			}

			ctx, span := startSpan(tm, m.Message)

			resp, err := handler(ctx, m.Message)

			if span != nil {
				endSpanWithError(span, err)
			}

			if err != nil {
				return err
			}

			err = s.Send(resp)
			if err != nil {
				return err
			}
		}
	}
}
