package tracing

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/tracing/metrics"
	"github.com/lolizeppelin/micro/tracing/tracers"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

func LoadMetricProvider(ctx context.Context, conf metrics.MetricConfig,
	res *resource.Resource) (*metric.MeterProvider, error) {
	loader := metrics.LoadExport(conf.Driver)
	if loader == nil {
		return nil, fmt.Errorf("metric export %s not found", conf.Driver)
	}
	exporter, err := loader(ctx, conf)
	if err != nil {
		return nil, err
	}

	var opts []metric.PeriodicReaderOption
	if conf.Batch.Timeout > 0 { // default 30s
		opts = append(opts, metric.WithTimeout(time.Duration(conf.Batch.Timeout)*time.Second))
	}
	if conf.Batch.Interval > 0 { // default 60s
		opts = append(opts, metric.WithInterval(time.Duration(conf.Batch.Interval)*time.Second))
	}
	return NewMetricProvider(ctx, metric.NewPeriodicReader(exporter), res)
}

func LoadTracerProvider(ctx context.Context, conf tracers.TracerConfig,
	res *resource.Resource) (*trace.TracerProvider, error) {
	loader := tracers.LoadExport(conf.Driver)
	if loader == nil {
		return nil, fmt.Errorf("tracer export %s not found", conf.Driver)
	}
	exporter, err := loader(ctx, conf)
	if err != nil {
		return nil, err
	}
	var opts []trace.BatchSpanProcessorOption
	if conf.Batch.Timeout > 0 {
		opts = append(opts, trace.WithBatchTimeout(time.Duration(conf.Batch.Timeout)*time.Second))
	}
	if conf.Batch.Size > 0 {
		opts = append(opts, trace.WithMaxExportBatchSize(conf.Batch.Size))
	}
	if conf.Batch.Queue > 0 {
		opts = append(opts, trace.WithMaxQueueSize(conf.Batch.Queue))
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, opts...),
		trace.WithResource(res),
	)
	return provider, nil
}
