package metrics

import (
	"context"
	"github.com/lolizeppelin/micro/utils/tls"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	exports = map[string]ExportMetric{}
)

func init() {
	exports["collector"] = NewGRPCExport
	exports["collector.grpc"] = NewGRPCExport
	exports["collector.http"] = NewHTTPExport
}

type MetricBatch struct {
	Timeout  int32 `json:"timeout,omitempty"` // 测试环境填1方便调试
	Interval int32 `json:"size,omitempty"`    // 测试环境填1方便调试
}

type MetricConfig struct {
	Driver      string                 `json:"driver" description:"Metric驱动"` // e.g prometheus/victoriametrics/greptime
	Endpoint    string                 `json:"endpoints" description:"Metric接口地址uri"`
	Batch       MetricBatch            `json:"batch" description:"批量上传配置"`
	Auth        map[string]string      `json:"auth,omitempty" description:"认证"`
	Options     map[string]any         `json:"options,omitempty" description:"驱动独立参数"`
	Credentials *tls.ClientCredentials `json:"credentials,omitempty" description:"ssl链接配置"`
}

type ExportMetric func(context.Context, MetricConfig) (metric.Exporter, error)

func LoadExport(driver string) ExportMetric {
	return exports[driver]
}
