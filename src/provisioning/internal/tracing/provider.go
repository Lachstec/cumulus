package tracing

import (
	"context"
	"errors"
	"github.com/Lachstec/mc-hosting/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// SetupOTEL initializes an OpenTelemetry Tracing Provider and registers it with the Telemtry SDK.
// The created provider pretty-prints to stdout by default.
func SetupOTEL(ctx context.Context, config config.TracingConfig) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error

		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	tracer, err := newTracer(ctx, config)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, tracer.Shutdown)
	otel.SetTracerProvider(tracer)

	return
}

func newTracer(ctx context.Context, config config.TracingConfig) (*trace.TracerProvider, error) {
	stdout, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	otlpClient := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(config.Endpoint), otlptracegrpc.WithInsecure())

	otlp, err := otlptrace.New(ctx, otlpClient)
	if err != nil {
		return nil, err
	}

	stdoutprocessor := trace.NewSimpleSpanProcessor(stdout)
	otlpprocessor := trace.NewSimpleSpanProcessor(otlp)

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String("1.0.0"),
	)

	provider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(stdoutprocessor),
		trace.WithSpanProcessor(otlpprocessor),
		trace.WithSampler(trace.AlwaysSample()),
	)

	return provider, nil
}
