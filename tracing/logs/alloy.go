package logs

func init() {

	/*
		https://grafana.org.cn/docs/alloy/latest/collect/opentelemetry-data/
	*/
	exports["alloy.http"] = NewHTTPExport

}
