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

// listar los modulos
func (h *ModuleHandler) GetModulesForTeacher(c *fiber.Ctx) error {

	claims := utils.GetClaims(c)

	// campos para paginar
	var paginated types.Paginated

	if err := c.QueryParser(&paginated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
			"data":    err,
		})
	}

	// validamos
	paginated.Validate()

	// obtenemos los modulos
	modules, details, err := data.GetModulesForTeacher(&paginated, claims.UserAPI.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	modulesApi := ModulesConverToAPI(modules, h.config.GetString("APP_HOST"))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    modulesApi,
		"details": details,
	})
}

// convierte las entidades de modulos a tipos de modulos para mostrar en la API REST xD
func ModulesConverToAPI(modules []data.Module, apphost string) []types.Module {
	// convertimos los modulos a types.modules
	modulesApi := make([]types.Module, len(modules))
	for i, module := range modules {
		modulesApi[i] = types.Module{
			ID:        module.ID,
			CreatedAt: module.CreatedAt.String(),
			UpdatedAt: module.UpdatedAt.String(),
			CreateBy: types.UserAPI{
				ID:        module.CreatedBy.ID,
				FirstName: module.CreatedBy.FirstName,
				LastName:  module.CreatedBy.LastName,
				Email:     module.CreatedBy.Email,
				URLAvatar: module.CreatedBy.URLAvatar,
			},
			Code:             module.Code,
			Title:            module.Title,
			ShortDescription: module.ShortDescription,
			TextRoot:         module.TextRoot,
			ImgBackURL:       apphost + module.ImgBackURL,
			Difficulty:       string(module.Difficulty),
			PointsToEarn:     module.PointsToEarn,
			Index:            module.Index,
			IsPublic:         module.IsPublic,
		}
	}
	return modulesApi
}

func (h *ModuleHandler) GetModules(c *fiber.Ctx) error {

	var paginated types.Paginated

	if err := c.QueryParser(&paginated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
			"data":    err,
		})
	}

	// validamos
	paginated.Validate()

	// obtenemos los modulos
	modules, details, err := data.GetModule(&paginated)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	modulesApi := ModulesConverToAPI(modules, h.config.GetString("APP_HOST"))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    modulesApi,
		"details": details,
	})
}

// un usuario se podra suscribir a un modulo
func (h *ModuleHandler) Subscribe(c *fiber.Ctx) error {

	claims := utils.GetClaims(c)

	var req types.ReqSubscription
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
			"data":    err,
		})
	}

	// validamos
	err := req.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	// creamos la suscripcion
	_, err = data.RegisterSubscription(claims.UserAPI.ID, req.Code)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})

}

// Subscriptions recupera todas las subscripciones de un usuario
func (h *ModuleHandler) Subscriptions(c *fiber.Ctx) error {

	var paginated types.Paginated

	if err := c.QueryParser(&paginated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
			"data":    err,
		})
	}

	// validamos
	paginated.Validate()

	claims := utils.GetClaims(c)

	// obtenemos los modulos
	modules, details, err := data.GetModuleForStudent(&paginated, claims.UserAPI.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	modulesApi := ModulesConverToAPI(modules, h.config.GetString("APP_HOST"))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    modulesApi,
		"details": details,
	})
}
