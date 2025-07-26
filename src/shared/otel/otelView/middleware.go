package otelView

import (
	"net/http"

	"giggler-golang/src/shared/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ApplyTracerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			instrumentedCtx, span := otel.Tracer().Start(
				r.Context(),
				"router",
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.KeyValue{
						Key:   "http.method",
						Value: attribute.StringValue(r.Method),
					},
					attribute.KeyValue{
						Key:   "http.url",
						Value: attribute.StringValue(r.URL.String()),
					},
				),
			)
			defer span.End()

			next.ServeHTTP(w, r.WithContext(instrumentedCtx))
		},
	)
}
