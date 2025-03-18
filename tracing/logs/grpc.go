package logs

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"net"
	"time"
)

/*
NewGRPCExport 创建一个使用 GRPC 协议连接的Exporter
*/
func NewGRPCExport(ctx context.Context, conf LoggingConfig) (log.Exporter, error) {
	endpoint := conf.Endpoint
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
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpointURL(endpoint),
		otlploggrpc.WithCompressor("gzip"),
		otlploggrpc.WithHeaders(conf.Auth),
		otlploggrpc.WithRetry(otlploggrpc.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlploggrpc.WithGRPCConn(conn),
		otlploggrpc.WithTimeout(5*time.Second))
}
