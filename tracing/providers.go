package tracing

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/tracing/logs"
	"github.com/lolizeppelin/micro/tracing/metrics"
	"github.com/lolizeppelin/micro/tracing/tracers"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

func LoadMetricProvider(ctx context.Context, conf metrics.MetricConfig,
	res *resource.Resource, producer ...metric.Producer) (*metric.MeterProvider, error) {
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
	if len(producer) > 0 {
		opts = append(opts, metric.WithProducer(producer[0]))
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)
	return provider, nil
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

func LoadLoggingProvider(ctx context.Context, conf logs.LoggingConfig,
	res *resource.Resource) (*log.LoggerProvider, error) {
	loader := logs.LoadExport(conf.Driver)
	if loader == nil {
		return nil, fmt.Errorf("metric export %s not found", conf.Driver)
	}
	exporter, err := loader(ctx, conf)
	if err != nil {
		return nil, err
	}

	var opts []log.BatchProcessorOption
	if conf.Batch.Timeout > 0 { // default 30s
		opts = append(opts, log.WithExportTimeout(time.Duration(conf.Batch.Timeout)*time.Second))
	}
	if conf.Batch.Interval > 0 { // default 60s
		opts = append(opts, log.WithExportInterval(time.Duration(conf.Batch.Interval)*time.Second))
	}

	provider := log.NewLoggerProvider(
		log.WithAttributeCountLimit(32),
		log.WithAttributeValueLengthLimit(1024),
		log.WithProcessor(log.NewBatchProcessor(exporter, opts...)),
		log.WithResource(res),
	)
	return provider, nil
}
