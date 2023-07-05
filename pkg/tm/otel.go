package tm

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/atuleu/otelog"
	"github.com/atuleu/otelog/pkg/hooks"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type otelProvider struct {
	shutdown func(context.Context) error
}

// Arguments needed for Open Telemetry
type OtelProviderArgs struct {
	// The Collector URL
	CollectorURL string
	// The Service Name
	ServiceName string
	// Its deployed version
	ServiceVersion string
	// Its verbose level for logs
	Level VerboseLevel
	// Force a Flush before Shutdown
	ForceFlushOnShutdown bool
}

type nullFormatter struct{}

func (f nullFormatter) Format(*logrus.Entry) ([]byte, error) {
	return nil, nil
}

func newOtelProvider(args OtelProviderArgs) Provider {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	conn, err := grpc.Dial(args.CollectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				MinConnectTimeout: 10 * time.Second,
				Backoff:           backoff.DefaultConfig,
			},
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second,
			Timeout: 20 * time.Second,
		}),
	)
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
		))
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(args.ServiceName),
		semconv.ServiceVersionKey.String(args.ServiceVersion),
		semconv.ServiceInstanceIDKey.String(hostname),
		semconv.HostIDKey.String(hostname),
	)

	spanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resources),
	))

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	metricExporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	metricReader := sdkmetric.NewPeriodicReader(metricExporter)
	otel.SetMeterProvider(sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resources),
		sdkmetric.WithReader(metricReader),
	))

	otel.SetTextMapPropagator(propagator)

	logExporter, err := otelog.NewLogExporter(
		otelog.WithResource(resources),
		otelog.WithScope(instrumentation.Scope{
			Name: "github.com/formicidae-tracker/olympus/pkg/tm",
		}),
		otelog.WithBatchLogProcessor(otelog.WithBatchTimeout(5*time.Second)),
		otelog.WithGRPCConn(conn),
	)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	otelog.SetLogExporter(logExporter)

	level := MapVerboseLevel(args.Level)
	hook := hooks.NewLogrusHook(hooks.FromLogrusLevel(level))

	logrus.SetLevel(level)
	logrus.AddHook(hook)

	var shutdown func(context.Context) error
	if args.ForceFlushOnShutdown == true {
		shutdown = func(ctx context.Context) error {
			return errors.Join(
				spanProcessor.ForceFlush(ctx),
				traceExporter.Shutdown(ctx),
				metricReader.ForceFlush(ctx),
				metricExporter.Shutdown(ctx),
			)
		}
	} else {
		shutdown = func(ctx context.Context) error {
			return errors.Join(
				traceExporter.Shutdown(ctx),
				metricExporter.Shutdown(ctx),
			)
		}
	}

	return &otelProvider{
		shutdown: shutdown,
	}
}

func (p *otelProvider) Shutdown(ctx context.Context) error {
	return p.shutdown(ctx)
}

func (p *otelProvider) NewLogger(domain string) *logrus.Entry {
	return logrus.WithField("domain", domain)
}

func (p *otelProvider) Enabled() bool {
	return true
}
