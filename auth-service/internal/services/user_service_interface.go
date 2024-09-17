package services

import (
	"context"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/requests"
	"github.com/saufiroja/go-otel/auth-service/internal/models"
)

type UserService interface {
	RegisterUser(ctx context.Context, request *requests.RegisterRequest) error
	LoginUser(ctx context.Context, request *requests.LoginRequest) (*models.Token, error)
}
