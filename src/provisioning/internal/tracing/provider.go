package tracing

import (
	"context"
	"github.com/Lachstec/mc-hosting/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}

func resources(config config.TracingConfig) *resource.Resource {
	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(config.ServiceName))
	return res
}

func NewTraceProvider(ctx context.Context, config config.TracingConfig) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(config.Endpoint),
		otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	res := resources(config)

	batchspan := sdktrace.NewBatchSpanProcessor(exporter)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(batchspan),
	)
	otel.SetTracerProvider(provider)

	propagator := newPropagator()
	otel.SetTextMapPropagator(propagator)

	return provider, nil
}
