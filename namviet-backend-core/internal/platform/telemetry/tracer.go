package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer initializes OpenTelemetry for distributed tracing.
// It tracks requests across Frontend, API Go, Redis, Supabase, and Gemini.
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// For production, you would configure an exporter here (e.g., OTLP to Jaeger/GCP Trace)
	// Example: exporter, err := otlptracegrpc.New(ctx)
	
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	log.Printf("OpenTelemetry Tracer initialized for service: %s", serviceName)
	return tp, nil
}
