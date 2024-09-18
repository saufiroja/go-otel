package repositories

import (
	"context"
	"github.com/saufiroja/go-otel/auth-service/internal/models"
	"github.com/saufiroja/go-otel/auth-service/pkg/databases"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := u.Trace.StartSpan(ctx, "repository.CreateUser")
	defer span.End()
	db := u.DB.Connection()

	query := `INSERT INTO users (user_id, full_name, email, password, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6)`
	span.SetAttributes(
		attribute.Key("user_id").String(user.UserId),
		attribute.Key("full_name").String(user.FullName),
		attribute.Key("email").String(user.Email),
		attribute.Key("created_at").Int64(user.CreatedAt.Unix()),
		attribute.Key("updated_at").Int64(user.UpdatedAt.Unix()),
	)

	span.AddEvent("executing SQL query", trace.WithAttributes(
		attribute.Key("sql.query").String(query),
	))

	_, err := db.ExecContext(ctx, query,
		user.UserId, user.FullName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		span.AddEvent("Failed to execute query", trace.WithAttributes(attribute.Key("error").String(err.Error())))
		span.SetStatus(codes.Error, "Error executing query")
		return err
	}

	span.AddEvent("Successfully created user", trace.WithAttributes(attribute.Key("userId").String(user.UserId)))
	span.SetStatus(codes.Ok, "Query executed successfully")

	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := u.Trace.StartSpan(ctx, "repository.GetUserByEmail")
	defer span.End()
	db := u.DB.Connection()

	query := `SELECT user_id, full_name, email, password, created_at, updated_at
				FROM users
				WHERE email = $1`

	row := db.QueryRowContext(ctx, query, email)

	span.SetAttributes(
		attribute.Key("email").String(email),
	)

	span.AddEvent("executing SQL query", trace.WithAttributes(
		attribute.Key("sql.query").String(query),
	))

	user := &models.User{}
	err := row.Scan(&user.UserId, &user.FullName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.AddEvent("user not found", trace.WithAttributes(attribute.Key("error").String(err.Error())))
		return nil, err
	}

	span.AddEvent("Successfully retrieved user data", trace.WithAttributes(attribute.Key("userId").String(user.UserId)))
	span.SetStatus(codes.Ok, "Query executed successfully")

	return user, nil
}
