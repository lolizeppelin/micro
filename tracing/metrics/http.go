package metrics

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"time"
)

/*
NewHTTPExport 创建一个使用 HTTP 协议连接的Exporter
*/
func NewHTTPExport(ctx context.Context, conf MetricConfig) (metric.Exporter, error) {
	tls, err := conf.Credentials.TLS()
	if err != nil {
		return nil, err
	}

	address := conf.Endpoint
	return otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpointURL(address),
		otlpmetrichttp.WithTLSClientConfig(tls),
		otlpmetrichttp.WithHeaders(conf.Auth),
		otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlpmetrichttp.WithTimeout(5*time.Second))
}
