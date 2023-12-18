package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var login types.Login
	if err := c.BodyParser(&login); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	user, ok, err := data.Login(login)

	if !ok || err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error al iniciar sesion", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": "tokensdkfjls",
		"user":  user,
	})
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	return c.SendString("Usuario registrado")
}
