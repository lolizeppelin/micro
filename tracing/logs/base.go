package logs

import (
	"context"
	"github.com/lolizeppelin/micro/utils/tls"
	"go.opentelemetry.io/otel/sdk/log"
)

var (
	exports = map[string]ExportLogging{}
)

func init() {
	exports["collector"] = NewGRPCExport
	exports["collector.grpc"] = NewGRPCExport
	exports["collector.http"] = NewHTTPExport
	exports["loki.http"] = NewHTTPExport
}

type LoggingBatch struct {
	Timeout  int32 `json:"timeout,omitempty"`
	Interval int32 `json:"interval,omitempty"`
}

type LoggingConfig struct {
	Driver      string                 `json:"driver" description:"Logging驱动"` // e.g  collector/alloy.http
	Endpoint    string                 `json:"endpoint" description:"Logging接口地址uri"`
	Disabled    bool                   `json:"disabled,omitempty" description:"是否禁用otel日志"`
	Batch       LoggingBatch           `json:"batch" description:"批量上传配置"`
	Auth        map[string]string      `json:"auth,omitempty" description:"认证"`
	Options     map[string]any         `json:"options,omitempty" description:"驱动独立参数"`
	Credentials *tls.ClientCredentials `json:"credentials,omitempty" description:"ssl链接配置"`
}

type ExportLogging func(context.Context, LoggingConfig) (log.Exporter, error)

func LoadExport(driver string) ExportLogging {
	return exports[driver]
}
