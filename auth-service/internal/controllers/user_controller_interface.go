package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
}
