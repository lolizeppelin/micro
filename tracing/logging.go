package tracing

import (
	"context"
	"github.com/lolizeppelin/micro/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"strings"
)

type LocalLOGExporter struct {
}

// ExportSpans 实现 trace.SpanExporter 接口
func (e *LocalLOGExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, span := range spans {
		log.Infof("**Exporting Span: %s (TraceID: %s, SpanID: %s)",
			span.Name(),
			span.SpanContext().TraceID(),
			span.SpanContext().SpanID(),
		)
		log.Infof("**  Metadata: %d %s", span.Status().Code, span.Status().Description)
		// 打印 Span 的其他属性
		for _, attr := range span.Attributes() {
			log.Infof("**    Attribute: %s = %v", attr.Key, attr.Value.AsInterface())
		}
		var errors []attribute.KeyValue
		// 打印 Span 事件
		for _, event := range span.Events() {
			log.Infof("**  Event: %s", event.Name)
			for _, attr := range event.Attributes {
				if attr.Key == "exception.message" || attr.Key == "exception.stacktrace" || event.Name == "http.error" {
					errors = append(errors, attr)
					continue
				}
				log.Infof("**      Event Attribute: %s = %s", attr.Key, attr.Value.AsString())
			}
		}

		for _, attr := range errors {
			if attr.Key == "exception.stacktrace" {
				log.Errorf("**    Stacktrace: %s", strings.TrimSpace(attr.Value.AsString()))
			} else if attr.Key == "exception.message" {
				log.Errorf("**  Errors: %s", strings.TrimSpace(attr.Value.AsString()))
			} else {
				log.Errorf("**  Errors key: %s | value: %s", attr.Key, strings.TrimSpace(attr.Value.AsString()))
			}
		}
	}
	return nil
}

// Shutdown 实现 trace.SpanExporter 接口
func (e *LocalLOGExporter) Shutdown(ctx context.Context) error {
	return nil
}

func NewLogExporter() *LocalLOGExporter {
	return &LocalLOGExporter{}
}

func NewLocalProvider(res *resource.Resource, fake ...bool) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter

	if len(fake) > 0 && fake[0] {
		exporter = &FakeExporter{}
	} else {
		exporter = NewLogExporter()
	}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(0),
			trace.WithMaxExportBatchSize(0),
		),
		trace.WithResource(res),
	)
	return provider, nil
}
