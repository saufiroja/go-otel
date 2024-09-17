package tracing

import (
	"context"
	"fmt"
	"github.com/saufiroja/go-otel/auth-service/config"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	Trace trace.Tracer
}

func NewExporter(ctx context.Context, logging logging.Logger, serviceName string, conf *config.AppConfig) *Tracer {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		))
	if err != nil {
		logging.LogError(fmt.Sprintf("failed to create resource: %v", err))
		panic(err)
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(conf.Otel.OTLPEndpoint),
	)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		logging.LogError(fmt.Sprintf("failed to create trace exporter: %v", err))
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	// set tracing global
	otel.SetTracerProvider(tp)

	logging.LogInfo("tracing initialized")

	return &Tracer{
		Trace: otel.Tracer(serviceName),
	}
}

func (t *Tracer) StartSpan(ctx context.Context, name string,
	opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.Trace.Start(ctx, name, opts...)
}
