package tm

import (
	"context"
	"os"
	"time"

	"github.com/atuleu/otelog"
	"github.com/atuleu/otelog/pkg/hooks"
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
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

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(args.CollectorURL),
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

	spanProcessor := sdktrace.NewBatchSpanProcessor(exporter)

	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resources),
	))

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	logExporter, err := otelog.NewLogExporter(
		otelog.WithResource(resources),
		otelog.WithScope(instrumentation.Scope{
			Name: "github.com/formicidae-tracker/olympus/pkg/tm",
		}),
		otelog.WithBatchLogProcessor(otelog.WithBatchTimeout(5*time.Second)),
		otelog.WithEndpoint(args.CollectorURL),
		otelog.WithInsecure(),
	)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	otelog.SetLogExporter(logExporter)

	level := MapVerboseLevel(args.Level)
	hook := hooks.NewLogrusHook(hooks.FromLogrusLevel(level))

	logrus.SetLevel(level)
	logrus.AddHook(hook)

	shutdown := exporter.Shutdown
	if args.ForceFlushOnShutdown == true {
		shutdown = func(ctx context.Context) error {
			spanProcessor.ForceFlush(ctx)
			return exporter.Shutdown(ctx)
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
