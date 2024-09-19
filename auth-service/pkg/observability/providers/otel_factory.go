package providers

import (
	"context"
	"fmt"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type ProviderFactory struct {
	logging logging.Logger
}

func NewProviderFactory(logging logging.Logger) *ProviderFactory {
	return &ProviderFactory{
		logging: logging,
	}
}

func (f *ProviderFactory) CreateResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		f.logging.LogError(fmt.Sprintf("failed to create resource: %v", err))
		return nil, err
	}
	return res, nil
}
