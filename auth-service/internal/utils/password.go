package utils

import (
	"context"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/tracing"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, hashedPassword, password string) error
}

type BcryptHasher struct {
	Trace *tracing.Tracer
}

func NewBcryptHasher(trace *tracing.Tracer) PasswordHasher {
	return &BcryptHasher{
		Trace: trace,
	}
}

func (b *BcryptHasher) Hash(ctx context.Context, password string) (string, error) {
	_, span := b.Trace.StartSpan(ctx, "utils.Hash")
	defer span.End()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (b *BcryptHasher) Compare(ctx context.Context, hashedPassword, password string) error {
	_, span := b.Trace.StartSpan(ctx, "utils.Compare")
	defer span.End()
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
