package otel

import (
	"context"
	"io"
	"runtime"
	"strings"
	"time"

	"giggler-golang/src/shared/errutil/must"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	logSDK "go.opentelemetry.io/otel/sdk/log"
	metricSDK "go.opentelemetry.io/otel/sdk/metric"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerVar trace.Tracer
	meterVar  metric.Meter
	loggerVar log.Logger
)

func init() {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	tracerExporter := must.Do(stdouttrace.New(stdouttrace.WithWriter(io.Discard)))
	meterExporter := must.Do(stdoutmetric.New(stdoutmetric.WithWriter(io.Discard)))
	loggerExporter := must.Do(stdoutlog.New(stdoutlog.WithWriter(io.Discard)))

	tracerProvider := traceSDK.NewTracerProvider(
		traceSDK.WithBatcher(
			tracerExporter,
			traceSDK.WithBatchTimeout(time.Second),
		),
	)
	meterProvider := metricSDK.NewMeterProvider(
		metricSDK.WithReader(
			metricSDK.NewPeriodicReader(meterExporter, metricSDK.WithInterval(3*time.Second)),
		),
	)
	logggerProvider := logSDK.NewLoggerProvider(
		logSDK.WithProcessor(
			logSDK.NewBatchProcessor(loggerExporter),
		),
	)

	tracerVar = tracerProvider.Tracer("giggler")
	meterVar = meterProvider.Meter("giggler")
	loggerVar = logggerProvider.Logger("giggler")
}

func Trace(ctx context.Context) (context.Context, trace.Span) {
	// Get the name of the calling function
	pc, _, _, _ := runtime.Caller(1)
	split := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	funcName := split[len(split)-1]

	return tracerVar.Start(ctx, funcName)
}

func Tracer() trace.Tracer {
	return tracerVar
}

func Meter(ctx context.Context) metric.Meter {
	return meterVar
}

func Logger() log.Logger {
	return loggerVar
}
