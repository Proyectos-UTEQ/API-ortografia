package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type QuestionHandler struct {
	config *viper.Viper
}

func NewQuestionHandler(config *viper.Viper) *QuestionHandler {
	return &QuestionHandler{
		config: config,
	}
}

// Registra una nueva pregunta para un modulo en concreto.
func (h *QuestionHandler) RegisterQuestionForModule(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id") // id del modulo
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   "Error al recuperar el id del modulo",
		})
	}
	iduint := uint(id)

	// Parseamos el body
	var question types.Question
	if err := c.BodyParser(&question); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validamos
	err = question.Validate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	question.ModuleID = &iduint

	// Registramos en la db.
	questionEntidad, err := data.RegisterQuestionForModule(question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(questionEntidad)
}

// Recupera todas las preguntas de un modulo
func (h *QuestionHandler) GetQuestionsForModule(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id") // id del modulo
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   "Error al recuperar el id del modulo",
		})
	}

	questions, err := data.GetQuestionsForModule(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": questions,
	})
}

func (h *QuestionHandler) DeleteQuestion(c *fiber.Ctx) error {
	idquestion, err := c.ParamsInt("idquestion")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   "Error al recuperar el id de la pregunta",
		})
	}

	err = data.DeleteQuestion(uint(idquestion))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *QuestionHandler) UpdateQuestion(c *fiber.Ctx) error {
	idquestion, err := c.ParamsInt("idquestion")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   "Error al recuperar el id de la pregunta",
		})
	}

	var question types.Question
	if err := c.BodyParser(&question); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	question.ID = uint(idquestion)

	err = data.UpdateQuestion(question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
