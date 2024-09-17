package repositories

import (
	"context"
	"github.com/saufiroja/go-otel/auth-service/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
