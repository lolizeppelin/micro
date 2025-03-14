package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	t_ver = oteltrace.WithInstrumentationVersion("1.0.0")
)

func GetTracer(name string) oteltrace.Tracer {
	provider := otel.GetTracerProvider()
	return provider.Tracer(name, t_ver)
}

func GetPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}

// Extract  通过字典中的OpenTelemetry数据,生成ctx(kafka接收时用,grpc已封装不需要主动Extract)
func Extract(ctx context.Context, carrier map[string]string) context.Context {
	propagator := GetPropagator()
	c := propagation.MapCarrier(carrier)
	ctx = propagator.Extract(ctx, c)
	return oteltrace.ContextWithRemoteSpanContext(context.Background(), oteltrace.SpanContextFromContext(ctx))
}

// Inject  OpenTelemetry数据注入字典中(kafka发送时用,grpc已封装不需要主动Inject)
func Inject(ctx context.Context) propagation.MapCarrier {
	propagator := GetPropagator()
	c := propagation.MapCarrier{}
	propagator.Inject(ctx, c)
	return c
}

func StartTrace(ctx context.Context, scope, name string,
	attributes ...attribute.KeyValue) (context.Context, oteltrace.Span) {
	tracer := GetTracer(scope)
	return tracer.Start(ctx, name,
		oteltrace.WithAttributes(attributes...),
	)
}

func FakeTraceProvider(res *resource.Resource) *trace.TracerProvider {
	exporter := &FakeExporter{}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(0),
			trace.WithMaxExportBatchSize(0),
		),
		trace.WithResource(res),
	)
	return provider
}

func LogTraceProvider(res *resource.Resource) *trace.TracerProvider {
	exporter := NewLogExporter()
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(0),
			trace.WithMaxExportBatchSize(0),
		),
		trace.WithResource(res),
	)
	return provider
}
