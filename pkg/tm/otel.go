package tm

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type otelProvider struct {
	mx       sync.Mutex
	spans    []trace.Span
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
}

type nullFormatter struct{}

func (f nullFormatter) Format(*logrus.Entry) ([]byte, error) {
	return nil, nil
}

func newOtelProvider(args OtelProviderArgs) Provider {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(args.CollectorURL),
		))
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	hostname, err := os.Hostname()
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

	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	))

	logrus.AddHook(otellogrus.NewHook(
		otellogrus.WithLevels(MapVerboseLevelList(args.Level)...),
	))
	logrus.SetFormatter(nullFormatter{})
	logrus.SetOutput(io.Discard)

	return &otelProvider{shutdown: exporter.Shutdown}
}

func (p *otelProvider) Shutdown(ctx context.Context) error {
	p.mx.Lock()
	defer p.mx.Unlock()

	for _, span := range p.spans {
		span.End()
	}
	return p.shutdown(ctx)
}

func (p *otelProvider) NewLogger(domain string) *logrus.Entry {
	p.mx.Lock()
	defer p.mx.Unlock()

	ctx, span := otel.Tracer("olympus").Start(context.Background(), domain)
	p.spans = append(p.spans, span)
	return logrus.WithContext(ctx)
}

func (p *otelProvider) Enabled() bool {
	return true
}
