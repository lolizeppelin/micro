package tracing

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	otelmetric "go.opentelemetry.io/otel/metric"
)

var MetricLevels = []logrus.Level{
	logrus.WarnLevel,
	logrus.ErrorLevel,
}

type LogrusMetricHook struct {
	errors   metric.Int64Counter
	warnings metric.Int64Counter
	options  otelmetric.MeasurementOption
}

func (h *LogrusMetricHook) Levels() []logrus.Level {
	return MetricLevels
}

func (h *LogrusMetricHook) Fire(entry *logrus.Entry) error {
	if entry.Level == logrus.WarnLevel {
		h.warnings.Add(context.Background(), 1, h.options)
	} else {
		h.errors.Add(context.Background(), 1, h.options)
	}
	return nil

}

type LogrusHook struct {
	Debug      bool
	attributes []otellog.KeyValue
	options    otelmetric.MeasurementOption

	errors   metric.Int64Counter
	warnings metric.Int64Counter

	logger otellog.Logger
}

func (h *LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LogrusHook) Fire(entry *logrus.Entry) error {
	var level otellog.Severity

	var event string
	ctx := entry.Context
	if ctx == nil {
		event = "log"
		ctx = context.Background()
	} else {
		event = "otel"
	}
	switch entry.Level {
	case logrus.FatalLevel:
		event = "fatal"
		level = otellog.SeverityError
		break
	case logrus.WarnLevel:
		h.warnings.Add(ctx, 1, h.options)
		level = otellog.SeverityWarn
		break
	case logrus.ErrorLevel:
		h.errors.Add(ctx, 1, h.options)
		level = otellog.SeverityError
		break
	case logrus.DebugLevel:
		level = otellog.SeverityDebug
		break
	case logrus.InfoLevel:
		level = otellog.SeverityInfo
		break
	case logrus.TraceLevel:
		level = otellog.SeverityTrace
		break
	}

	record := otellog.Record{}
	record.SetEventName(event)
	record.SetSeverity(level)
	record.SetTimestamp(entry.Time)
	record.SetBody(otellog.StringValue(entry.Message))

	//attrs := h.attributes
	//for k, v := range entry.Data {
	//	attrs = append(attrs, otellog.String(k, fmt.Sprintf("%v", v)))
	//}
	record.AddAttributes(h.attributes...)
	h.logger.Emit(entry.Context, record)
	return nil
}

func NewLogrusHook(attributes []attribute.KeyValue, attrs []otellog.KeyValue, logger ...otellog.Logger) logrus.Hook {
	meter := GetMeter(ScopeName, _version)
	ec, _ := meter.Int64Counter("log.err")
	wc, _ := meter.Int64Counter("log.warn")

	if len(logger) > 0 {
		return &LogrusHook{
			logger:     logger[0],
			errors:     ec,
			warnings:   wc,
			attributes: attrs,
			options:    otelmetric.WithAttributes(attributes...),
		}
	}

	return &LogrusMetricHook{
		errors:   ec,
		warnings: wc,
		options:  otelmetric.WithAttributes(attributes...),
	}

}
