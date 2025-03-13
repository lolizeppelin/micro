package tracing

import "C"
import (
	"context"
	"github.com/lolizeppelin/micro/utils/tls"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OTELProvider interface {
	ForceFlush(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

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

type MetricBatch struct {
	Timeout  int32 `json:"timeout,omitempty"` // 测试环境填1方便调试
	Interval int32 `json:"size,omitempty"`    // 测试环境填1方便调试
}

type MetricConfig struct {
	Driver      string                 `json:"driver" description:"Metric驱动"` // e.g prometheus/victoriametrics/greptime
	Endpoints   string                 `json:"endpoints" description:"Metric接口地址uri"`
	Batch       MetricBatch            `json:"batch" description:"批量上传配置"`
	Auth        map[string]string      `json:"auth,omitempty" description:"认证"`
	Options     map[string]any         `json:"options,omitempty" description:"驱动独立参数"`
	Credentials *tls.ClientCredentials `json:"credentials,omitempty" description:"ssl链接配置"`
}

type OTELConf struct {
	Disabled bool          `json:"disabled" description:"禁用追踪"`
	Jaeger   *TracerConfig `json:"tracer" description:"Tracer配置"`
	Metric   *MetricConfig `json:"metric" description:"Metric配置"`
}

type FakeExporter struct {
}

func (e *FakeExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

func (e *FakeExporter) Shutdown(ctx context.Context) error {
	return nil
}
