package metrics

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"time"
)

/*
NewGRPCExport 创建一个使用 GRPC 协议连接的Exporter
*/
func NewGRPCExport(ctx context.Context, conf MetricConfig) (metric.Exporter, error) {
	endpoint := conf.Endpoint
	creds, err := conf.Credentials.Credentials()
	if err != nil {
		return nil, err
	}
	//uri, err := url.Parse(endpoint)
	//if err != nil {
	//	return nil, fmt.Errorf("parse OTEL metrics endpoint failed")
	//}
	//if !utils.IncludeInSlice([]string{"http", "https"}, uri.Scheme) {
	//	return nil, fmt.Errorf("OTEL metrics endpoint scheme %s not supported", uri.Scheme)
	//}
	//options := []grpc.DialOption{
	//	grpc.WithTransportCredentials(creds),
	//	grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
	//		return net.DialTimeout("tcp", uri.Host, 5)
	//	}),
	//}
	//conn, err := grpc.NewClient(uri.Host, options...)
	//if err != nil {
	//	return nil, err
	//}
	return otlpmetricgrpc.New(ctx,
		//otlpmetricgrpc.WithGRPCConn(conn),
		otlpmetricgrpc.WithEndpointURL(endpoint),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithHeaders(conf.Auth),
		otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlpmetricgrpc.WithTLSCredentials(creds),
		otlpmetricgrpc.WithTimeout(5*time.Second))
}
