package metrics

func init() {

	/*
		https://docs.victoriametrics.com/guides/getting-started-with-opentelemetry/app.go-collector.example
	*/
	exports["victoria"] = NewGRPCExport
	exports["victoriametrics"] = NewGRPCExport
	exports["VictoriaMetrics"] = NewGRPCExport

}
