package tm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/atuleu/otelog"
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
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
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

func connectToCollector(endpoint string) (*grpc.ClientConn, error) {
	return grpc.Dial(endpoint,
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
}

func buildTraceExporter(conn *grpc.ClientConn, resource *resource.Resource) (*trace.TracerProvider, trace.SpanProcessor, error) {

	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
		))

	if err != nil {
		return nil, nil, err
	}

	spanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resource),
	)

	return traceProvider, spanProcessor, nil
}

func buildMetricExporter(conn *grpc.ClientConn, resource *resource.Resource) (*metric.MeterProvider, sdkmetric.Reader, error) {

	metricExporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, nil, err
	}

	metricReader := sdkmetric.NewPeriodicReader(metricExporter)
	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(metricReader),
	)

	return metricProvider, metricReader, nil
}

func newOtelProvider(args OtelProviderArgs) Provider {

	fatal := func(err error) {
		slog.With("error", err).Error("could not initialize Open Telemetry")
		os.Exit(127)
	}

	fatal(fmt.Errorf("telemetry is not yet implemented with log/slog package"))

	hostname, err := os.Hostname()
	if err != nil {
		fatal(err)
	}

	conn, err := connectToCollector(args.CollectorURL)
	if err != nil {
		fatal(err)
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNamespaceKey.String("formicidae-tracker"),
		semconv.ServiceNameKey.String(args.ServiceName),
		semconv.ServiceVersionKey.String(args.ServiceVersion),
		semconv.ServiceInstanceIDKey.String(hostname),
		semconv.HostIDKey.String(hostname),
	)

	traceProvider, spanProcessor, err := buildTraceExporter(conn, resource)
	if err != nil {
		fatal(err)
	}
	otel.SetTracerProvider(traceProvider)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	meterProvider, metricReader, err := buildMetricExporter(conn, resource)
	if err != nil {
		fatal(err)
	}
	otel.SetMeterProvider(meterProvider)

	logExporter, err := otelog.NewLogExporter(
		otelog.WithResource(resource),
		otelog.WithScope(instrumentation.Scope{
			Name: "github.com/formicidae-tracker/olympus/pkg/tm",
		}),
		otelog.WithBatchLogProcessor(otelog.WithBatchTimeout(5*time.Second)),
		otelog.WithGRPCConn(conn),
	)
	if err != nil {
		fatal(err)
	}
	otelog.SetLogExporter(logExporter)

	var shutdown func(context.Context) error
	if args.ForceFlushOnShutdown == true {
		shutdown = func(ctx context.Context) error {
			return errors.Join(
				spanProcessor.ForceFlush(ctx),
				traceProvider.Shutdown(ctx),
				metricReader.ForceFlush(ctx),
				meterProvider.Shutdown(ctx),
			)
		}
	} else {
		shutdown = func(ctx context.Context) error {
			return errors.Join(
				traceProvider.Shutdown(ctx),
				meterProvider.Shutdown(ctx),
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

func (p *otelProvider) NewLogger(domain string) *slog.Logger {
	return logger.With("domain", domain)
}

func (p *otelProvider) Enabled() bool {
	return true
}
