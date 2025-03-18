package tracing

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/tracing/logs"
	"github.com/lolizeppelin/micro/tracing/metrics"
	"github.com/lolizeppelin/micro/tracing/tracers"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	ScopeName = "micro/tracing"
)

var (
	_version, _ = micro.NewVersion("1.0.0")
)

type OTELProvider interface {
	ForceFlush(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type OTELConf struct {
	Disabled bool                  `json:"disabled" description:"禁用追踪"`
	Tracer   *tracers.TracerConfig `json:"tracer" description:"Tracer配置"`
	Metric   *metrics.MetricConfig `json:"metric" description:"Metric配置"`
	Logging  *logs.LoggingConfig   `json:"logging" description:"Logging配置"`
}

type FakeExporter struct {
}

func (e *FakeExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

func (e *FakeExporter) Shutdown(ctx context.Context) error {
	return nil
}
