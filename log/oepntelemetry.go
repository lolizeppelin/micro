package log

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OTELExporter struct {
	skip bool
}

// ExportSpans 实现 trace.SpanExporter 接口
func (e *OTELExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {

	if e.skip {
		return nil
	}

	for _, span := range spans {
		LOG.Infof("**Exporting Span: %s (TraceID: %s, SpanID: %s)\n",
			span.Name(),
			span.SpanContext().TraceID(),
			span.SpanContext().SpanID(),
		)

		LOG.Infof("**  Metadata: %d %s\n", span.Status().Code, span.Status().Description)
		// 打印 Span 的其他属性
		for _, attr := range span.Attributes() {
			LOG.Infof("**    Attribute: %s = %v\n", attr.Key, attr.Value.AsInterface())
		}
		var errors []attribute.KeyValue
		// 打印 Span 事件
		for _, event := range span.Events() {
			LOG.Infof("**  Event: %s\n", event.Name)
			for _, attr := range event.Attributes {
				errors = append(errors, attr)
				if attr.Key == "exception.message" || attr.Key == "exception.stacktrace" {
					errors = append(errors, attr)
				}
				LOG.Infof("**      Event Attribute: %s = %v\n", attr.Key, attr.Value.AsInterface())
			}
		}

		for _, attr := range errors {
			if attr.Key == "exception.message" {
				LOG.Errorf("**  Errors: %s\n", attr.Value.AsString())
			} else {
				LOG.Errorf("**    Stacktrace: %s\n", attr.Value.AsString())
			}
		}
	}
	return nil
}

// Shutdown 实现 trace.SpanExporter 接口
func (e *OTELExporter) Shutdown(ctx context.Context) error {
	return nil
}

func NewOTELExporter(skip bool) *OTELExporter {
	return &OTELExporter{
		skip: skip,
	}
}
