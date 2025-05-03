package tracing

import (
	"context"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"log/slog"
)

func LoadSlogHandler(name string, provider *log.LoggerProvider) slog.Handler {
	return otelslog.NewHandler(name, otelslog.WithLoggerProvider(provider),
		otelslog.WithVersion(_version.Version(true)))
}

func LoadSlogMeterHandler() slog.Handler {

	meter := GetMeter(ScopeName, _version)
	ec, _ := meter.Int64Counter("log.err")
	wc, _ := meter.Int64Counter("log.warn")

	return &SlogMeterHandler{
		error:   ec,
		warning: wc,
	}

}

type SlogMeterHandler struct {
	error   metric.Int64Counter
	warning metric.Int64Counter
}

func (h *SlogMeterHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level == slog.LevelWarn || level == slog.LevelError
}

func (h *SlogMeterHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level == slog.LevelWarn {
		h.warning.Add(ctx, 1)
	} else if record.Level == slog.LevelError {
		h.error.Add(ctx, 1)
	}
	return nil
}

func (h *SlogMeterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *SlogMeterHandler) WithGroup(name string) slog.Handler {
	return h
}
