package tracers

func init() {

	/*
		火山云 https://www.volcengine.com/docs/6470/812322
		gRPC 协议的请求头（Header）中配置鉴权参数：
		map[string]string{
			"x-tls-otel-tracetopic": topicId,
			"x-tls-otel-ak": ak,
			"x-tls-otel-sk": sk,
			"x-tls-otel-region": region}
	*/
	exports["volcengine"] = NewGRPCExport

}
