package tracers

func init() {

	exports["jaeger"] = NewGRPCExport
	exports["Jaeger"] = NewGRPCExport

	exports["jaeger.http"] = NewHTTPExport

}
