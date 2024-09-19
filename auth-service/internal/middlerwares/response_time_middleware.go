package middlerwares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/metrics"
	"time"
)

type Middleware struct {
	Meter *metrics.Metric
}

func NewMiddleware(meter *metrics.Metric) *Middleware {
	return &Middleware{
		Meter: meter,
	}
}

func (m *Middleware) ResponseTimeMiddleware(ctx context.Context, apiName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		if err := c.Next(); err != nil {
			return err
		}

		duration := time.Since(startTime)

		m.Meter.Histogram(ctx, apiName, "Response time", "ms", duration)

		return nil
	}
}
