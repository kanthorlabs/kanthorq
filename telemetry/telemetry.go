package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var once sync.Once
var Telemetry *OpenTelemetry

func init() {
	once.Do(func() {
		Telemetry = &OpenTelemetry{}
	})
}

func Propagator() propagation.TextMapPropagator {
	return Telemetry.Propagator()
}

func Tracer() trace.Tracer {
	return Telemetry.Tracer()
}

func Meter() metric.Meter {
	return Telemetry.Meter()
}

func MeterCounter(name string, options ...metric.Int64CounterOption) func(incr int64, options ...metric.AddOption) {
	counter, err := Telemetry.Meter().Int64Counter(name, options...)
	if err != nil {
		return func(incr int64, options ...metric.AddOption) {}
	}

	return func(incr int64, options ...metric.AddOption) {
		counter.Add(context.Background(), incr, options...)
	}
}

func MeterHistogram(name string, options ...metric.Float64HistogramOption) func(incr float64, options ...metric.RecordOption) {
	options = append(options, metric.WithUnit("s"))
	histogram, err := Telemetry.Meter().Float64Histogram(name, options...)
	if err != nil {
		return func(incr float64, options ...metric.RecordOption) {}
	}

	return func(incr float64, options ...metric.RecordOption) {
		histogram.Record(context.Background(), incr, options...)
	}
}
