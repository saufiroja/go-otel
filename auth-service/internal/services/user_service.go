package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/requests"
	"github.com/saufiroja/go-otel/auth-service/internal/models"
	"github.com/saufiroja/go-otel/auth-service/internal/repositories"
	"github.com/saufiroja/go-otel/auth-service/internal/utils"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"github.com/saufiroja/go-otel/auth-service/pkg/tracing"
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
		u.Logger.LogError(fmt.Sprintf("User with email: %s already exists", request.Email))
		return fmt.Errorf("user with email: %s already exists", request.Email)
	}

	// hash password
	password, err := u.PasswordHasher.Hash(ctx, request.Password)
	if err != nil {
		u.Logger.LogError(fmt.Sprintf("Error hashing password: %v", err))
		return fmt.Errorf("error hashing password: %v", err)
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

	return u.UserRepository.CreateUser(ctx, userModel)
}

func (u *userService) LoginUser(ctx context.Context, request *requests.LoginRequest) (*models.Token, error) {
	ctx, span := u.Trace.StartSpan(ctx, "service.LoginUser")
	defer span.End()

	u.Logger.LogInfo(fmt.Sprintf("login user with email %s", request.Email))

	user, err := u.UserRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		u.Logger.LogError(fmt.Sprintf("Error getting user by email: %v", err))
		return nil, fmt.Errorf("error getting user by email: %v", err)
	}

	err = u.PasswordHasher.Compare(ctx, user.Password, request.Password)
	if err != nil {
		u.Logger.LogError(fmt.Sprintf("Error comparing password: %v", err))
		return nil, fmt.Errorf("error comparing password: %v", err)
	}

	// Generate token
	token := &requests.GenerateTokenRequest{
		UserId:   user.UserId,
		FullName: user.FullName,
	}

	accessToken, _, err := u.GenerateToken.GenerateAccessToken(ctx, token)
	if err != nil {
		u.Logger.LogError(fmt.Sprintf("Error generating access token: %v", err))
		return nil, fmt.Errorf("error generating access token: %v", err)
	}

	refreshToken, _, err := u.GenerateToken.GenerateRefreshToken(ctx, token)
	if err != nil {
		u.Logger.LogError(fmt.Sprintf("Error generating refresh token: %v", err))
		return nil, fmt.Errorf("error generating refresh token: %v", err)
	}

	res := &models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return res, nil
}
