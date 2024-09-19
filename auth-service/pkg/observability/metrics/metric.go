package metrics

import (
	"context"
	"fmt"
	"github.com/saufiroja/go-otel/auth-service/config"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/providers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkMetrics "go.opentelemetry.io/otel/sdk/metric"
	"time"
)

type Metric struct {
	Meter       metric.Meter
	logger      logging.Logger
	ServiceName string
}

func NewMetric(ctx context.Context, serviceName string, provider *providers.ProviderFactory,
	conf *config.AppConfig, logging logging.Logger) *Metric {
	res, err := provider.CreateResource(ctx, serviceName)
	metricsExp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(conf.Otel.OTLPEndpoint),
	)
	if err != nil {
		logging.LogPanic(fmt.Sprintf("failed to create OTLP metric exporter: %v", err))
	}

	meterProvider := sdkMetrics.NewMeterProvider(
		sdkMetrics.WithResource(res),
		sdkMetrics.WithReader(sdkMetrics.NewPeriodicReader(
			metricsExp)),
	)

	otel.SetMeterProvider(meterProvider)

	logging.LogInfo("metrics initialized")

	return &Metric{
		Meter:       otel.Meter(serviceName),
		logger:      logging,
		ServiceName: serviceName,
	}
}

func (m *Metric) Counter(ctx context.Context, name string, description, unit string) {
	counter, err := m.Meter.Int64Counter(fmt.Sprintf("%s.%s", m.ServiceName, name),
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		m.logger.LogError(fmt.Sprintf("failed to create counter: %v", err))
		return
	}
	counter.Add(ctx, 1)
}

func (m *Metric) Histogram(ctx context.Context, name string, description, unit string, duration time.Duration) {
	histogram, err := m.Meter.Int64Histogram(fmt.Sprintf("%s.%s", m.ServiceName, name),
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		m.logger.LogError(fmt.Sprintf("failed to create histogram: %v", err))
		return
	}

	histogram.Record(ctx, duration.Milliseconds())
}
