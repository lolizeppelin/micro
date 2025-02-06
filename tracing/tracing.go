package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	version = oteltrace.WithInstrumentationVersion("1.0.0")
)

func GetTracer(name string) oteltrace.Tracer {
	provider := otel.GetTracerProvider()
	return provider.Tracer(name, version)
}

func GetPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}
