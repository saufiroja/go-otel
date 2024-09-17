package repositories

import (
	"context"
	"github.com/saufiroja/go-otel/auth-service/internal/models"
	"github.com/saufiroja/go-otel/auth-service/pkg/databases"
	"github.com/saufiroja/go-otel/auth-service/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

type userRepository struct {
	DB    databases.PostgresManager
	Trace *tracing.Tracer
}

func NewUserRepository(db databases.PostgresManager, trace *tracing.Tracer) UserRepository {
	return &userRepository{
		DB:    db,
		Trace: trace,
	}
}

func (u *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	ctx, span := u.Trace.StartSpan(ctx, "repository.CreateUser", trace.WithAttributes())
	defer span.End()
	db := u.DB.Connection()

	query := `INSERT INTO users (user_id, full_name, email, password, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExecContext(ctx, query,
		user.UserId, user.FullName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := u.Trace.StartSpan(ctx, "repository.GetUserByEmail", trace.WithAttributes())
	defer span.End()
	db := u.DB.Connection()

	query := `SELECT user_id, full_name, email, password, created_at, updated_at
				FROM users
				WHERE email = $1`

	row := db.QueryRowContext(ctx, query, email)

	user := &models.User{}
	err := row.Scan(&user.UserId, &user.FullName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}
