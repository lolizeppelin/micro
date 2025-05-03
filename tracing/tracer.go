package tracing

import (
	"context"
	"github.com/lolizeppelin/micro/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"strings"
)

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

type TracerLog struct {
}

// ExportSpans 实现 trace.SpanExporter 接口
func (e *TracerLog) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, span := range spans {
		log.Infof(ctx, "**Exporting Span: %s (TraceID: %s, SpanID: %s)",
			span.Name(),
			span.SpanContext().TraceID(),
			span.SpanContext().SpanID(),
		)
		log.Infof(ctx, "**  Metadata: %d %s", span.Status().Code, span.Status().Description)
		// 打印 Span 的其他属性
		for _, attr := range span.Attributes() {
			log.Infof(ctx, "**    Attribute: %s = %v", attr.Key, attr.Value.AsInterface())
		}
		var errors []attribute.KeyValue
		// 打印 Span 事件
		for _, event := range span.Events() {
			log.Infof(ctx, "**  Event: %s", event.Name)
			for _, attr := range event.Attributes {
				if attr.Key == "exception.message" || attr.Key == "exception.stacktrace" || event.Name == "http.error" {
					errors = append(errors, attr)
					continue
				}
				log.Infof(ctx, "**      Event Attribute: %s = %s", attr.Key, attr.Value.AsString())
			}
		}

		for _, attr := range errors {
			if attr.Key == "exception.stacktrace" {
				log.Errorf(ctx, "**    Stacktrace: %s", strings.TrimSpace(attr.Value.AsString()))
			} else if attr.Key == "exception.message" {
				log.Errorf(ctx, "**  Errors: %s", strings.TrimSpace(attr.Value.AsString()))
			} else {
				log.Errorf(ctx, "**  Errors key: %s | value: %s", attr.Key, strings.TrimSpace(attr.Value.AsString()))
			}
		}
	}
	return nil
}

// Shutdown 实现 trace.SpanExporter 接口
func (e *TracerLog) Shutdown(ctx context.Context) error {
	return nil
}

func NewLogExporter() *TracerLog {
	return &TracerLog{}
}
