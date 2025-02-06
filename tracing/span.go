package tracing

import (
	"context"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type propagateKey struct{}

var _PropagateCtxKey = propagateKey{}

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
	propagator := GetPropagator()
	propagator.Inject(ctx, propagation.MapCarrier(carrier))
	return context.WithValue(ctx, _PropagateCtxKey, carrier)
}
