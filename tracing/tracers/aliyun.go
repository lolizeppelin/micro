package tracers

func init() {

	/*
		https://help.aliyun.com/zh/opentelemetry/user-guide/use-managed-service-for-opentelemetry-to-submit-the-trace-data-of-a-go-application?spm=a2c4g.11186623.help-menu-90275.d_2_0_4_0.1c452299deKr1t
	*/
	exports["aliyun"] = NewGRPCExport

}
