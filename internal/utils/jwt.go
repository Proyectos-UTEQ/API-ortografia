package utils

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/gofiber/fiber/v2"
)

func GetClaims(c *fiber.Ctx) *types.UserClaims {
	return c.Locals("user").(*types.UserClaims)
}
