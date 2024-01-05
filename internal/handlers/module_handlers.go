package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"log"

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
	// recuperamos los claims del usuarios
	claims := utils.GetClaims(c)

	// Parseamos el body
	var module types.Module
	if err := c.BodyParser(&module); err != nil {
		log.Println("Error al registrar modulo", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
		})
	}

	// Validar datos para registro inicial.
	resp, err := types.Validate(&module)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error en la validacion de datos",
			"data":    resp,
		})
	}

	// Crea el modulo en la base de datos y recuperamos los datos del usuario
	// que creo el modulo.
	moduleResponse, err := data.RegisterModuleForTeacher(&module, claims.UserAPI.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al registrar modulo",
			"data":    err,
		})
	}

	// Generamos la url de la imagen del modulo.
	moduleResponse.ImgBackURL = h.config.GetString("APP_HOST") + moduleResponse.ImgBackURL

	return c.Status(fiber.StatusCreated).JSON(moduleResponse)

}
