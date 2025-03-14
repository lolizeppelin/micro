package metrics

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"net"
	"time"
)

/*
NewGRPCExport 创建一个使用 GRPC 协议连接的Exporter
*/
func NewGRPCExport(ctx context.Context, conf MetricConfig) (metric.Exporter, error) {
	endpoint := conf.Endpoints
	creds, err := conf.Credentials.Credentials()
	if err != nil {
		return nil, err
	}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return net.DialTimeout("tcp", endpoint, 5)
		}),
	}
	conn, err := grpc.NewClient(endpoint, options...)
	if err != nil {
		return nil, err
	}
	return otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpointURL(endpoint),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithHeaders(conf.Auth),
		otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlpmetricgrpc.WithGRPCConn(conn),
		otlpmetricgrpc.WithTimeout(5*time.Second))
}
