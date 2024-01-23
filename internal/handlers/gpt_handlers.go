package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/interfaces"
	"Proyectos-UTEQ/api-ortografia/internal/services"
	"Proyectos-UTEQ/api-ortografia/pkg/types"

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

func (g *GPTHandler) GenerateQuestion(c *fiber.Ctx) error {
	var reqContext types.GPT

	if err := c.BodyParser(&reqContext); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// result, err := services.NewGPT(g.config.GetString("APP_OPENAI_API_KEY")).GPTTest(reqContext.Context)

	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	if err := reqContext.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var serviceIA interfaces.QuestionIA
	switch reqContext.Server {
	case "chatgptapi":
		serviceIA = services.NewServicoIA(g.config)
	case "gpts4u":
		serviceIA = services.NewServiceNewGpts4u(g.config)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	result, err := serviceIA.GenerateQuestion(reqContext.TypeQuestion, reqContext.Context)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": result,
	})
}
