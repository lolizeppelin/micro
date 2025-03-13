/**
标准单位https://ucum.org/ucum
*/

package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	m_ver = otelmetric.WithInstrumentationVersion("1.0.0")
)

func GetMeter(name string) otelmetric.Meter {
	provider := otel.GetMeterProvider()
	return provider.Meter(name, m_ver)
}

/*
GetCounter 单调递增计数器，比如可以用来记录订单数、总的请求数。
*/
func GetCounter(scope, name string) otelmetric.Int64Counter {
	return GetIntCounter(scope, name, "", "")
}

/*
GetIntCounter 单调递增计数器，比如可以用来记录订单数、总的请求数。
*/
func GetIntCounter(scope, name, unit, desc string) otelmetric.Int64Counter {
	counter, _ := GetMeter(scope).Int64Counter(
		name,
		otelmetric.WithUnit(unit),
		otelmetric.WithDescription(desc),
	)

	return counter
}

func GetFloatCounter(scope, name, unit, desc string) otelmetric.Float64Counter {
	counter, _ := GetMeter(scope).Float64Counter(
		name,
		otelmetric.WithUnit(unit),
		otelmetric.WithDescription(desc),
	)
	return counter
}

/*
GetUpDownCounter 可以减少的计数器
*/
func GetUpDownCounter(scope, name, unit, desc string) otelmetric.Int64UpDownCounter {
	counter, _ := GetMeter(scope).Int64UpDownCounter(
		name,
		otelmetric.WithUnit(unit),
		otelmetric.WithDescription(desc),
	)
	return counter
}

/*
GetTimeHistogram 通常用于记录请求延迟、响应时间等
*/
func GetTimeHistogram(scope, name string) otelmetric.Float64Histogram {

	histogram, _ := GetMeter(scope).Float64Histogram(
		name,
		otelmetric.WithUnit("second"),
		otelmetric.WithExplicitBucketBoundaries(0.5, 1, 3, 5, 10, 15, 30),
	)
	return histogram
}

/*
GetIntHistogram 通常用于记录请求延迟、响应时间等
*/
func GetIntHistogram(scope, name, unit, desc string) otelmetric.Int64Histogram {
	histogram, _ := GetMeter(scope).Int64Histogram(
		name,
		otelmetric.WithUnit(unit),
		otelmetric.WithDescription(desc),
	)
	return histogram
}

/*
GetGauge 随时变化的值, cpu、内存等
*/
func GetGauge(scope, name string) otelmetric.Int64Gauge {
	return GetIntGauge(scope, name, "", "")
}

/*
GetIntGauge 随时变化的值, cpu、内存等
*/
func GetIntGauge(scope, name, unit, desc string) otelmetric.Int64Gauge {
	gauge, _ := GetMeter(scope).Int64Gauge(
		name,
		otelmetric.WithUnit(unit),
		otelmetric.WithDescription(desc),
	)
	return gauge
}

func NewMetricProvider(ctx context.Context, exporter metric.Reader, res *resource.Resource) (*metric.MeterProvider, error) {
	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)
	return provider, nil
}
