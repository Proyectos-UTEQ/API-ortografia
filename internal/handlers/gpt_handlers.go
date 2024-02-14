package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/services"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type GPTHandler struct {
	config *viper.Viper
}

func NewGPTHandler(config *viper.Viper) *GPTHandler {
	return &GPTHandler{
		config: config,
	}
}

// GenerateResponse Servicio para que el frontend pueda generar una respuesta para su ChatGPT 3.5 elegante.
func (g *GPTHandler) GenerateResponse(c *fiber.Ctx) error {
	var req types.RequestGPT
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	gpt := services.NewGPT(g.config)
	res, err := gpt.GenerateResponse(req.Request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	return c.JSON(types.ResponseGPT{
		Response: res,
	})
}

func (g *GPTHandler) GenerateImage(c *fiber.Ctx) error {
	var req types.RequestGPT
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ia := services.NewGPT(g.config)

	url, err := ia.GenerateImage(req.Request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}

func (g *GPTHandler) GenerateQuestion(c *fiber.Ctx) error {
	var reqContext types.GPT

	// Parseamos los datos que nos envía el frontend
	if err := c.BodyParser(&reqContext); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := reqContext.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ia := services.NewGPT(g.config)
	result, err := ia.GenerateQuestion(reqContext.TypeQuestion, reqContext.Context)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": result,
	})
}
