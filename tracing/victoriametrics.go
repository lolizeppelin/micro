package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"
)

/*
NewVictoriaProvider 创建一个使用 HTTP 协议连接的VictoriaMetrics Exporter
https://docs.victoriametrics.com/guides/getting-started-with-opentelemetry/app.go-collector.example
*/
func NewVictoriaProvider(ctx context.Context, conf MetricConfig, res *resource.Resource) (*metric.MeterProvider, error) {
	//
	if conf.Driver != "victoria" {
		return nil, fmt.Errorf("metric dirver error")
	}
	tls, err := conf.Credentials.TLS()
	if err != nil {
		return nil, err
	}

	address := conf.Endpoints
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpointURL(address),
		otlpmetrichttp.WithTLSClientConfig(tls),
		otlpmetrichttp.WithHeaders(conf.Auth),
		otlpmetrichttp.WithTemporalitySelector(
			func(kind metric.InstrumentKind) metricdata.Temporality {
				return metricdata.DeltaTemporality
			},
		),
		otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlpmetrichttp.WithTimeout(5*time.Second))

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
