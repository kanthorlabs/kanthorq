package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var once sync.Once

func init() {
	once.Do(func() {
		exporter, err := otlptracegrpc.New(context.Background())
		if err != nil {
			panic(err)
		}

		var tracer = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.Default()),
		)
		otel.SetTracerProvider(tracer)
	})
}

var Tracer = otel.Tracer("github.com/kanthorlabs/kanthorq")
