package tracers

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"net"
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

	address := conf.Endpoint

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

	var exporter *otlptrace.Exporter
	exporter, err = otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(address),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{ // 重试机制
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
		}),
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	var opts []trace.BatchSpanProcessorOption
	if conf.Batch.Timeout > 0 {
		opts = append(opts, trace.WithBatchTimeout(time.Duration(conf.Batch.Timeout)*time.Second))
	}
	if conf.Batch.Size > 0 {
		opts = append(opts, trace.WithMaxExportBatchSize(conf.Batch.Size))
	}
	if conf.Batch.Queue > 0 {
		opts = append(opts, trace.WithMaxQueueSize(conf.Batch.Queue))
	}

	return exporter, nil
}
