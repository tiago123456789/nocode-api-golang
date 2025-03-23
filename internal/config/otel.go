package config

import (
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewTraceExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func NewMetricExpoter() (metric.Exporter, error) {
	return stdoutmetric.New()
}

func NewPropagation() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)
}

func NewTraceProvider(traceExporter trace.SpanExporter) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			traceExporter,
			trace.WithBatchTimeout(time.Second),
		),
	)

	return traceProvider
}

func NewMeterProvider(meterExporter metric.Exporter) *metric.MeterProvider {
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				meterExporter,
				metric.WithInterval(time.Second*10),
			),
		),
	)

	return meterProvider
}
