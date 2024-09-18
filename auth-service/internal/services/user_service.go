package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/requests"
	"github.com/saufiroja/go-otel/auth-service/internal/models"
	"github.com/saufiroja/go-otel/auth-service/internal/repositories"
	"github.com/saufiroja/go-otel/auth-service/internal/utils"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"time"
)

type userService struct {
	UserRepository repositories.UserRepository
	Logger         logging.Logger
	GenerateToken  *utils.GenerateToken
	Trace          *tracing.Tracer
	PasswordHasher utils.PasswordHasher
}

func NewUserService(userRepository repositories.UserRepository, logger logging.Logger,
	generateToken *utils.GenerateToken, trace *tracing.Tracer, passwordHasher utils.PasswordHasher) UserService {
	return &userService{
		UserRepository: userRepository,
		Logger:         logger,
		GenerateToken:  generateToken,
		Trace:          trace,
		PasswordHasher: passwordHasher,
	}
}

func (u *userService) RegisterUser(ctx context.Context, request *requests.RegisterRequest) error {
	ctx, span := u.Trace.StartSpan(ctx, "service.RegisterUser")
	defer span.End()

	u.Logger.LogInfo(fmt.Sprintf(
		"Registering user with email %s, full name %s", request.Email, request.FullName))

	// Check if user already exists
	_, err := u.UserRepository.GetUserByEmail(ctx, request.Email)
	if err == nil {
		span.SetAttributes(attribute.Key("error.email").String(request.Email))
		span.AddEvent("User already exists")
		span.SetStatus(codes.Error, "User already exists")
		u.Logger.LogError(fmt.Sprintf("User with email: %s already exists", request.Email))
		return errors.New("user already exists")
	}

	span.SetAttributes(
		attribute.Key("email").String(request.Email),
		attribute.Key("full_name").String(request.FullName),
		attribute.Key("password").String(request.Password),
	)

	// hash password
	password, err := u.PasswordHasher.Hash(ctx, request.Password)
	if err != nil {
		span.AddEvent("Failed to hash password")
		span.SetStatus(codes.Error, "Error hashing password")
		u.Logger.LogError(fmt.Sprintf("Error hashing password: %v", err))
		return errors.New("error hashing password")
	}

	// Create user
	userModel := &models.User{
		UserId:    uuid.New().String(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = u.UserRepository.CreateUser(ctx, userModel)
	if err != nil {
		span.AddEvent("Failed to create user")
		span.SetStatus(codes.Error, "Error creating user")
		u.Logger.LogError(fmt.Sprintf("Error creating user: %v", err))
		return errors.New("error creating user")
	}

	span.AddEvent("User created successfully")
	span.SetStatus(codes.Ok, "User created successfully")

	return nil
}

func (u *userService) LoginUser(ctx context.Context, request *requests.LoginRequest) (*models.Token, error) {
	ctx, span := u.Trace.StartSpan(ctx, "service.LoginUser")
	defer span.End()

	u.Logger.LogInfo(fmt.Sprintf("login user with email %s", request.Email))

	user, err := u.UserRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		span.SetAttributes(attribute.Key("error.email").String(request.Email))
		span.AddEvent("Failed to get user by email")
		span.SetStatus(codes.Error, "Error getting user by email")
		u.Logger.LogError(fmt.Sprintf("Error getting user by email: %v", err))
		return nil, errors.New("email not found")
	}

	span.SetAttributes(attribute.Key("email").String(user.Email))
	span.SetAttributes(attribute.Key("user_id").String(user.UserId))
	span.SetAttributes(attribute.Key("full_name").String(user.FullName))

	err = u.PasswordHasher.Compare(ctx, user.Password, request.Password)
	if err != nil {
		span.AddEvent("password mismatch")
		span.SetStatus(codes.Error, "Password mismatch")
		u.Logger.LogError(fmt.Sprintf("password mismatch: %v", err))
		return nil, errors.New("password mismatch")
	}

	// Generate token
	token := &requests.GenerateTokenRequest{
		UserId:   user.UserId,
		FullName: user.FullName,
	}

	accessToken, _, err := u.GenerateToken.GenerateAccessToken(ctx, token)
	if err != nil {
		span.AddEvent("Failed to generate access token")
		span.SetStatus(codes.Error, "Error generating access token")
		u.Logger.LogError(fmt.Sprintf("Error generating access token: %v", err))
		return nil, errors.New("error generating access token")
	}

	refreshToken, _, err := u.GenerateToken.GenerateRefreshToken(ctx, token)
	if err != nil {
		span.AddEvent("Failed to generate refresh token")
		span.SetStatus(codes.Error, "Error generating refresh token")
		u.Logger.LogError(fmt.Sprintf("Error generating refresh token: %v", err))
		return nil, errors.New("error generating refresh token")
	}

	res := &models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	span.AddEvent("Login successful")
	span.SetStatus(codes.Ok, "Login successful")

	return res, nil
}
