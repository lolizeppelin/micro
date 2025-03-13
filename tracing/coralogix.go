package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"net"
	"time"
)

/*
NewCoralogixProvider 创建一个使用 GRPC 协议连接的Coralogix Exporter
https://coralogix.com/docs/opentelemetry/instrumentation-options/golang-opentelemetry-instrumentation/
*/
func NewCoralogixProvider(ctx context.Context, conf MetricConfig, res *resource.Resource) (*metric.MeterProvider, error) {

	if conf.Driver != "coralogix" {
		return nil, fmt.Errorf("metric dirver error")
	}
	cred, err := conf.Credentials.Credentials()
	if err != nil {
		return nil, err
	}

	address := conf.Endpoints
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return net.DialTimeout("tcp", address, 5)
		}),
	}

	conn, err := grpc.NewClient(address, options...)
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpointURL(address),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithHeaders(conf.Auth),
		otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlpmetricgrpc.WithGRPCConn(conn),
		otlpmetricgrpc.WithTimeout(5*time.Second))

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
