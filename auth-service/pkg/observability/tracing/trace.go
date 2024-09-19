package tracing

import (
	"context"
	"fmt"
	"github.com/saufiroja/go-otel/auth-service/config"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/providers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is a wrapper for OpenTelemetry Tracer
type Tracer struct {
	Trace trace.Tracer
}

// NewTracer creates a new tracer provider
func NewTracer(ctx context.Context, serviceName string, provider *providers.ProviderFactory,
	conf *config.AppConfig, logging logging.Logger) *Tracer {
	res, err := provider.CreateResource(ctx, serviceName)
	if err != nil {
		logging.LogPanic(err.Error())
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(conf.Otel.OTLPEndpoint),
	)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		logging.LogPanic(fmt.Sprintf("failed to create OTLP trace exporter: %v", err))
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

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.Trace.Start(ctx, name, opts...)
}
