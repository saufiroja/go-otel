package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/go-otel/auth-service/config"
	"github.com/saufiroja/go-otel/auth-service/internal/controllers"
	"github.com/saufiroja/go-otel/auth-service/internal/middlerwares"
	"github.com/saufiroja/go-otel/auth-service/internal/repositories"
	"github.com/saufiroja/go-otel/auth-service/internal/services"
	"github.com/saufiroja/go-otel/auth-service/internal/utils"
	"github.com/saufiroja/go-otel/auth-service/pkg/databases"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/metrics"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/providers"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/tracing"
)

type App struct {
	*fiber.App
}

func NewApp() *App {
	return &App{
		App: fiber.New(),
	}
}

func (a *App) Start() {
	logger := logging.NewLogrusAdapter()
	conf := config.NewAppConfig(logger)
	postgresInstance := databases.NewPostgres(conf, logger)
	defer postgresInstance.CloseConnection()

	const serviceName = "auth-service"
	resource := providers.NewProviderFactory(logger)
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx, serviceName, resource, conf, logger)
	meter := metrics.NewMetric(ctx, serviceName, resource, conf, logger)

	//middlewares
	responseTimeMiddleware := middlerwares.NewMiddleware(meter)
	//utils
	generateToken := utils.NewGenerateToken(conf, tracer)
	passwordHasher := utils.NewBcryptHasher(tracer)

	userRepository := repositories.NewUserRepository(postgresInstance, tracer)
	userService := services.NewUserService(userRepository, logger, generateToken, tracer, passwordHasher)
	userController := controllers.NewUserController(userService, tracer, meter)

	a.Post("/register",
		responseTimeMiddleware.ResponseTimeMiddleware(ctx, "register"), userController.RegisterUser)

	a.Post("/login",
		responseTimeMiddleware.ResponseTimeMiddleware(ctx, "login"), userController.LoginUser)

	if err := a.Listen(fmt.Sprintf(":%s", conf.Http.Port)); err != nil {
		logger.LogPanic(err.Error())
	}
}
