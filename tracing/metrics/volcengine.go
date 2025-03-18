package metrics

func init() {

	/*
		火山云 https://www.volcengine.com/docs/6431/97735
		gRPC 协议的请求头（Header）中配置鉴权参数：
		map[string]string{
			"X-ByteAPM-AppKey": app_key
		}

		export OTEL_GO_X_RESOURCE=true
		export OTEL_EXPORTER_OTLP_PROTOCOL=<protocal>
		export OTEL_EXPORTER_OTLP_ENDPOINT=<apmplus_endpoint>
		export OTEL_SERVICE_NAME=<service_name>
		export OTEL_EXPORTER_OTLP_HEADERS="X-ByteAPM-AppKey=<app_key>"

	*/
	exports["volcengine"] = NewGRPCExport
	//exports["volcengine"] = NewHTTPExport

}
