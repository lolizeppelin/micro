package tracers

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/url"
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
	endpoint := conf.Endpoint
	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse jaeper endpoint failed")
	}
	if !utils.IncludeInSlice([]string{"http", "https"}, uri.Scheme) {
		return nil, fmt.Errorf("OTEL jaeper endpoint scheme %s not supported", uri.Scheme)
	}
	return otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(uri.Host),
		otlptracehttp.WithTLSClientConfig(tls),
		otlptracehttp.WithHeaders(conf.Auth),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}), otlptracehttp.WithTimeout(5*time.Second))

}
