package metrics

func init() {

	/*
	   https://docs.greptime.com/zh/greptimecloud/integrations/otlp
	*/
	exports["greptime"] = NewGRPCExport
	exports["GreptimeDB"] = NewGRPCExport

}
