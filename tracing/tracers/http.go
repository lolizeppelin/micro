package tracers

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

/*
NewHTTPExport 创建一个使用 HTTP 协议连接的Exporter
*/
func NewHTTPExport(ctx context.Context, conf TracerConfig) (trace.SpanExporter, error) {

	tls, err := conf.Credentials.TLS()
	if err != nil {
		return nil, err
	}
	address := conf.Endpoint
	return otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(address),
		otlptracehttp.WithTLSClientConfig(tls),
		otlptracehttp.WithHeaders(conf.Auth),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}), otlptracehttp.WithTimeout(5*time.Second))

}
