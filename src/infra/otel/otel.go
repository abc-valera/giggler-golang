package otel

// TODO: finish

// import (
// 	"io"
// 	"time"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
// 	"go.opentelemetry.io/otel/log"
// 	"go.opentelemetry.io/otel/metric"
// 	"go.opentelemetry.io/otel/propagation"
// 	logSDK "go.opentelemetry.io/otel/sdk/log"
// 	metricSDK "go.opentelemetry.io/otel/sdk/metric"
// 	traceSDK "go.opentelemetry.io/otel/sdk/trace"
// 	"go.opentelemetry.io/otel/trace"
// )

// func Init() (trace.Tracer, metric.Meter, log.Logger, error) {
// 	propagator := propagation.NewCompositeTextMapPropagator(
// 		propagation.TraceContext{},
// 		propagation.Baggage{},
// 	)
// 	otel.SetTextMapPropagator(propagator)

// 	tracerExporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	tracerProvider := traceSDK.NewTracerProvider(
// 		traceSDK.WithBatcher(
// 			tracerExporter,
// 			traceSDK.WithBatchTimeout(time.Second),
// 		),
// 	)
// 	tracer := tracerProvider.Tracer("giggler")

// 	metricExporter, err := stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	metricProvider := metricSDK.NewMeterProvider(
// 		metricSDK.WithReader(
// 			metricSDK.NewPeriodicReader(metricExporter, metricSDK.WithInterval(3*time.Second)),
// 		),
// 	)
// 	metric := metricProvider.Meter("giggler")

// 	loggerExporter, err := stdoutlog.New(stdoutlog.WithWriter(io.Discard))
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	loggerProvider := logSDK.NewLoggerProvider(
// 		logSDK.WithProcessor(
// 			logSDK.NewBatchProcessor(loggerExporter),
// 		),
// 	)
// 	logger := loggerProvider.Logger("giggler")

// 	return tracer, metric, logger, nil
// }
