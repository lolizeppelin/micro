package tracers

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

/*
NewGRPCExport 创建一个使用 GRPC 协议连接的Exporter
*/
func NewGRPCExport(ctx context.Context, conf TracerConfig) (trace.SpanExporter, error) {
	cred, err := conf.Credentials.Credentials()
	if err != nil {
		return nil, err
	}
	endpoint := conf.Endpoint
	//uri, err := url.Parse(endpoint)
	//if err != nil {
	//	return nil, fmt.Errorf("parse OTEL tracer endpoint failed")
	//}
	//if !utils.IncludeInSlice([]string{"http", "https"}, uri.Scheme) {
	//	return nil, fmt.Errorf("OTEL tracer endpoint scheme %s not supported", uri.Scheme)
	//}
	//options := []grpc.DialOption{
	//	grpc.WithTransportCredentials(cred),
	//	grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
	//		return net.DialTimeout("tcp", uri.Host, 5)
	//	}),
	//}
	//conn, err := grpc.NewClient(uri.Host, options...)
	//if err != nil {
	//	return nil, err
	//}

	return otlptracegrpc.New(ctx,
		//otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlptracegrpc.WithTLSCredentials(cred),
		otlptracegrpc.WithTimeout(5*time.Second))
}
