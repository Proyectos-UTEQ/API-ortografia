package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	data.Login()

	return c.SendString("hola, mundo")
}
