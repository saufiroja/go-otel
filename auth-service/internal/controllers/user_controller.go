package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/requests"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/responses"
	"github.com/saufiroja/go-otel/auth-service/internal/services"
	"github.com/saufiroja/go-otel/auth-service/pkg/observability/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type userController struct {
	UserService services.UserService
	Trace       *tracing.Tracer
}

func NewUserController(userService services.UserService, trace *tracing.Tracer) UserController {
	return &userController{
		UserService: userService,
		Trace:       trace,
	}
}

func (u *userController) RegisterUser(c *fiber.Ctx) error {
	ctx, span := u.Trace.StartSpan(c.Context(), "controller.RegisterUser")
	defer span.End()

	request := &requests.RegisterRequest{}
	err := c.BodyParser(request)
	if err != nil {
		span.SetAttributes(
			attribute.Key("error.email").String(request.Email),
			attribute.Key("error.full_name").String(request.FullName),
			attribute.Key("error.password").String(request.Password),
		)
		span.AddEvent("Failed to parse request body")
		span.SetStatus(codes.Error, "Bad request body")
		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	span.SetAttributes(
		attribute.Key("email").String(request.Email),
		attribute.Key("full_name").String(request.FullName),
		attribute.Key("password").String(request.Password),
	)

	err = u.UserService.RegisterUser(ctx, request)
	if err != nil {
		span.AddEvent("Failed to register user",
			trace.WithAttributes(attribute.Key("error.email").String(request.Email)))
		span.SetStatus(codes.Error, err.Error())
		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	span.AddEvent("User registered successfully")
	span.SetStatus(codes.Ok, "User registered successfully")

	responseSuccess := responses.NewResponse[any](
		"User registered successfully", fiber.StatusCreated, nil)
	return c.Status(fiber.StatusCreated).JSON(responseSuccess)
}

func (u *userController) LoginUser(c *fiber.Ctx) error {
	ctx, span := u.Trace.StartSpan(c.Context(), "controller.LoginUser")
	defer span.End()

	request := &requests.LoginRequest{}
	err := c.BodyParser(request)
	if err != nil {
		span.SetAttributes(attribute.Key("error.email").String(request.Email))
		span.AddEvent("Failed to parse request body")
		span.SetStatus(codes.Error, "Bad request body")

		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	span.SetAttributes(attribute.Key("email").String(request.Email))

	token, err := u.UserService.LoginUser(ctx, request)
	if err != nil {
		span.AddEvent("Login failed for user",
			trace.WithAttributes(attribute.Key("error.email").String(request.Email)))
		span.SetStatus(codes.Error, err.Error())

		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	span.AddEvent("User logged in successfully")
	span.SetStatus(codes.Ok, "User logged in successfully")

	responseSuccess := responses.NewResponse[any](
		"User logged in successfully", fiber.StatusOK, token)
	return c.Status(fiber.StatusOK).JSON(responseSuccess)
}
