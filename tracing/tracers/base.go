package tracers

import (
	"context"
	"github.com/lolizeppelin/micro/utils/tls"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	exports = map[string]ExportTracer{}
)

type TracerBatch struct {
	Timeout int32 `json:"timeout,omitempty"` // 测试环境填1方便调试
	Size    int   `json:"size,omitempty"`
	Queue   int   `json:"queue,omitempty"`
}

type TracerConfig struct {
	Driver      string                 `json:"driver" description:"Tracer驱动"` // e.g jaeger
	Endpoint    string                 `json:"endpoints" description:"接口地址"`
	Batch       TracerBatch            `json:"batch" description:"批量上传配置"`
	Auth        map[string]string      `json:"auth,omitempty" description:"认证"`
	Options     map[string]any         `json:"options,omitempty" description:"驱动独立参数参数"`
	Credentials *tls.ClientCredentials `json:"credentials,omitempty" description:"ssl链接配置"`
}

type ExportTracer func(context.Context, TracerConfig) (trace.SpanExporter, error)

func LoadExport(driver string) ExportTracer {
	return nil
}
