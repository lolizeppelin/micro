package tracing

import (
	"context"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type propagateKey struct{}

var _PropagateCtxKey = propagateKey{}

// Extract  from carrier into context
func Extract(carrier map[string]string) context.Context {
	propagator := GetPropagator()
	c := propagation.MapCarrier(carrier)
	ctx := propagator.Extract(context.Background(), c)
	return context.WithValue(ctx, _PropagateCtxKey, c)
}

// Inject  from ctx into carrier
func Inject(ctx context.Context, carrier map[string]string) context.Context {
	propagator := GetPropagator()
	c := propagation.MapCarrier(carrier)
	propagator.Inject(ctx, c)
	return context.WithValue(ctx, _PropagateCtxKey, c)
}

func ExtractSpan(ctx context.Context) oteltrace.SpanContext {
	span := oteltrace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext()
	}
	carrier, ok := ctx.Value(_PropagateCtxKey).(propagation.MapCarrier)
	if !ok {
		return oteltrace.SpanContext{}
	}
	propagator := GetPropagator()
	extractedCtx := propagator.Extract(ctx, carrier)
	extractedSpan := oteltrace.SpanFromContext(extractedCtx)
	return extractedSpan.SpanContext()
}

func InjectSpan(ctx context.Context, carrier map[string]string) context.Context {
	span := oteltrace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ctx
	}
	return Inject(ctx, carrier)
}
