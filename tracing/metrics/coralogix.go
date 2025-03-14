package metrics

func init() {

	/*
		https://coralogix.com/docs/opentelemetry/instrumentation-options/golang-opentelemetry-instrumentation/
	*/
	exports["coralogix"] = NewGRPCExport
	exports["Coralogix"] = NewGRPCExport

}
