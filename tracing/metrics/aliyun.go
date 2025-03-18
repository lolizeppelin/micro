package metrics

func init() {

	/*
		https://help.aliyun.com/zh/arms/tracing-analysis/use-opentelemetry-to-submit-trace-data-of-go-applications
	*/
	exports["aliyun"] = NewGRPCExport

}
