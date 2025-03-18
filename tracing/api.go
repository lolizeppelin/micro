/**
标准单位https://ucum.org/ucum
*/

package tracing

import (
	"context"
	"github.com/lolizeppelin/micro"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func GetMeter(name string, version *micro.Version) otelmetric.Meter {
	provider := otel.GetMeterProvider()
	return provider.Meter(name, otelmetric.WithInstrumentationVersion(version.Version(true)))
}

func GetTracer(name string, version *micro.Version) oteltrace.Tracer {
	provider := otel.GetTracerProvider()
	return provider.Tracer(name, oteltrace.WithInstrumentationVersion(version.Version(true)))
}

//func GetLogger(name string, version *micro.Version) otellog.Logger {
//	provider := otellog.NewLoggerConfig()
//	return provider.Tracer(name, oteltrace.WithInstrumentationVersion(version.Version(true)))
//}

func StartTrace(ctx context.Context, scope, name string, version *micro.Version,
	attributes ...attribute.KeyValue) (context.Context, oteltrace.Span) {
	tracer := GetTracer(scope, version)
	return tracer.Start(ctx, name,
		oteltrace.WithAttributes(attributes...),
	)
}
