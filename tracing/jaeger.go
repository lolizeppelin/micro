package tracing

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type JaegerBatch struct {
	Timeout time.Duration `json:"timeout,omitempty"` // 测试环境填1方便调试
	Size    int           `json:"size,omitempty"`
	Queue   int           `json:"queue,omitempty"`
}

type JaegerConfig struct {
	Endpoint string      `json:"endpoint"`
	Batch    JaegerBatch `json:"batch"`
}

func NewJaegerProvider(ctx context.Context, conf JaegerConfig, res *resource.Resource) (*trace.TracerProvider, error) {
	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(conf.Endpoint),
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}
	var opts []trace.BatchSpanProcessorOption
	if conf.Batch.Timeout > 0 {
		opts = append(opts, trace.WithBatchTimeout(conf.Batch.Timeout))
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
