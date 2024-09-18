package utils

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saufiroja/go-otel/auth-service/config"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/requests"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/tracing"
	"time"
)

type GenerateToken struct {
	Secret string
	Trace  *tracing.Tracer
}

func NewGenerateToken(conf *config.AppConfig, trace *tracing.Tracer) *GenerateToken {
	return &GenerateToken{
		Secret: conf.Jwt.Secret,
		Trace:  trace,
	}
}

func (g *GenerateToken) GenerateAccessToken(ctx context.Context,
	request *requests.GenerateTokenRequest) (string, int64, error) {
	ctx, span := g.Trace.StartSpan(ctx, "utils.GenerateToken.GenerateAccessToken")
	defer span.End()

	expired := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   request.UserId,
		"full_name": request.FullName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(g.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expired, nil
}

func (g *GenerateToken) GenerateRefreshToken(ctx context.Context,
	request *requests.GenerateTokenRequest) (string, int64, error) {
	ctx, span := g.Trace.StartSpan(ctx, "utils.GenerateToken.GenerateRefreshToken")
	defer span.End()

	expired := time.Now().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   request.UserId,
		"full_name": request.FullName,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(g.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expired, nil
}
