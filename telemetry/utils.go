package telemetry

import (
	"fmt"

	"go.opentelemetry.io/otel/propagation"
)

func MapCarrier(kv map[string]any) propagation.TextMapCarrier {
	carrier := propagation.MapCarrier{}
	for k, v := range kv {
		if k == "traceparent" {
			carrier.Set(k, fmt.Sprintf("%v", v))
		}
	}
	return carrier
}
