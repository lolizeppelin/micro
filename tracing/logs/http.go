package logs

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"time"
)

/*
NewHTTPExport 创建一个使用 HTTP 协议连接的Exporter
*/
func NewHTTPExport(ctx context.Context, conf LoggingConfig) (otellog.Exporter, error) {
	tls, err := conf.Credentials.TLS()
	if err != nil {
		return nil, err
	}
	address := conf.Endpoint
	return otlploghttp.New(ctx,
		otlploghttp.WithEndpointURL(address),
		otlploghttp.WithTLSClientConfig(tls),
		otlploghttp.WithHeaders(conf.Auth),
		otlploghttp.WithRetry(otlploghttp.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlploghttp.WithTimeout(5*time.Second))
}
