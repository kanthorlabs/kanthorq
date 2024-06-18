package telemetry

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

type OpenTelemetry struct {
	propagator propagation.TextMapPropagator
	meter      metric.MeterProvider
	tracer     trace.TracerProvider

	cmu sync.Mutex
	mmu sync.Mutex
	tmu sync.Mutex
}

func (tel *OpenTelemetry) Propagator() propagation.TextMapPropagator {
	if tel.propagator == nil {
		tel.propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
	return tel.propagator
}

func (tel *OpenTelemetry) Meter() metric.Meter {
	tel.mmu.Lock()
	defer tel.mmu.Unlock()

	if tel.meter == nil {
		tel.meter = noopmetric.NewMeterProvider()
	}
	return tel.meter.Meter("github.com/kanthorlabs/kanthorq")
}

func (tel *OpenTelemetry) Tracer() trace.Tracer {
	tel.tmu.Lock()
	defer tel.tmu.Unlock()

	if tel.tracer == nil {
		tel.tracer = nooptrace.NewTracerProvider()
	}
	return tel.tracer.Tracer("github.com/kanthorlabs/kanthorq")
}

func (tel *OpenTelemetry) Start(ctx context.Context) error {
	tel.cmu.Lock()
	defer tel.cmu.Unlock()

	meterExp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return err
	}
	tel.meter = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(meterExp)),
		sdkmetric.WithResource(resource.Default()),
	)

	traceExp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return err
	}
	tel.tracer = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(resource.Default()),
	)

	return nil
}

func (tel *OpenTelemetry) Stop(ctx context.Context) error {
	tel.cmu.Lock()
	defer tel.cmu.Unlock()

	var merr error
	if tel.meter != nil {
		if meter, ok := tel.meter.(*sdkmetric.MeterProvider); ok {
			if err := meter.Shutdown(ctx); err != nil {
				merr = errors.Join(merr, err)
			}
		}
	}
	if tel.tracer != nil {
		if tracer, ok := tel.tracer.(*sdktrace.TracerProvider); ok {
			if err := tracer.Shutdown(ctx); err != nil {
				merr = errors.Join(merr, err)
			}
		}
	}
	return nil
}
