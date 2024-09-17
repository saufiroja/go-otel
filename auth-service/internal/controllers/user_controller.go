package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/requests"
	"github.com/saufiroja/go-otel/auth-service/internal/contracts/responses"
	"github.com/saufiroja/go-otel/auth-service/internal/services"
	"github.com/saufiroja/go-otel/auth-service/pkg/tracing"
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
		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = u.UserService.RegisterUser(ctx, request)
	if err != nil {
		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

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
		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	token, err := u.UserService.LoginUser(ctx, request)
	if err != nil {
		response := responses.NewResponse[any](
			err.Error(), fiber.StatusBadRequest, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	responseSuccess := responses.NewResponse[any](
		"User logged in successfully", fiber.StatusOK, token)
	return c.Status(fiber.StatusOK).JSON(responseSuccess)
}
