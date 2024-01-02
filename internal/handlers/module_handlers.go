package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type ModuleHandler struct {
	config *viper.Viper
}

// NewModuleHandler crea un nuevo handler de modulos.
func NewModuleHandler(config *viper.Viper) *ModuleHandler {
	return &ModuleHandler{
		config: config,
	}
}

func (h *ModuleHandler) CreateModuleForTeacher(c *fiber.Ctx) error {
	claims := utils.GetClaims(c)
	return c.SendString(fmt.Sprintf("Modulo creado para el tipo de usuario %s con nombre %s", claims.TypeUser, claims.FirstName))
}
