package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"fmt"
	"log"
	"strconv"

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

// UpdateModule Actualiza el modulo en la base de datos.
func (h *ModuleHandler) UpdateModule(c *fiber.Ctx) error {

	// claims := utils.GetClaims(c)

	idModule := c.Params("id")
	if idModule == "" {
		log.Println("Error al registrar modulo")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al registrar modulo",
		})
	}

	// convertir a uint el id module
	id, err := strconv.Atoi(idModule)
	if err != nil {
		log.Println("Error al registrar modulo", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al registrar modulo",
		})
	}

	var module types.Module
	// Parseamos el body
	if err := c.BodyParser(&module); err != nil {
		log.Println("Error al registrar modulo", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
		})
	}

	// seteamos el id del modulo.
	module.ID = uint(id)

	// Validar datos.
	resp, err := types.Validate(&module)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error en la validacion de datos",
			"data":    resp,
		})
	}

	// Actualizamos el modulo en la db
	moduleData, err := data.UpdateModule(&module)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al registrar modulo",
			"data":    err,
		})
	}

	moduleResponse := data.ModuleToApi(*moduleData)

	return c.Status(fiber.StatusOK).JSON(moduleResponse)
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
	_ = paginated.Validate()

	// obtenemos los modulos
	modules, details, err := data.GetModulesForTeacher(&paginated, claims.UserAPI.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	modulesApi := data.ModulesToAPI(modules, h.config.GetString("APP_HOST"))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    modulesApi,
		"details": details,
	})
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
	_ = paginated.Validate()

	// obtenemos los modulos
	modules, details, err := data.GetModule(&paginated)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	modulesApi := data.ModulesToAPI(modules, h.config.GetString("APP_HOST"))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    modulesApi,
		"details": details,
	})
}

func (h *ModuleHandler) GetModuleWithIsSubscribed(c *fiber.Ctx) error {

	claims := utils.GetClaims(c)

	var paginated types.Paginated

	if err := c.QueryParser(&paginated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear los datos",
			"data":    err,
		})
	}

	// validamos
	_ = paginated.Validate()

	// obtenemos los modulos
	modules, details, err := data.GetModuleWithUserSubscription(&paginated, claims.UserAPI.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    err,
		})
	}

	modulesApi := data.ModuleUserSubcriptionToApi(modules)

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
	_ = paginated.Validate()

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

	modulesApi := data.ModulesToAPI(modules, h.config.GetString("APP_HOST"))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    modulesApi,
		"details": details,
	})
}

func (h *ModuleHandler) GetStudents(c *fiber.Ctx) error {
	idModule := c.Params("id")

	// Convertir el idModule a uint
	id, err := strconv.Atoi(idModule)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear el id del modulo",
		})
	}

	// Obtener los estudiantes del modulo
	students, err := data.GetStudentsByModule(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	studentsData := data.UsersToAPI(students)

	return c.Status(fiber.StatusOK).JSON(studentsData)
}

func (h *ModuleHandler) GetModuleByID(c *fiber.Ctx) error {
	idModule := c.Params("id")

	// Convertir el idModule a uint
	id, err := strconv.Atoi(idModule)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear el id del modulo",
		})
	}

	module, err := data.ModuleByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	moduleResponse := data.ModuleToApi(*module)

	return c.JSON(moduleResponse)
}

func (h *ModuleHandler) GenerateTest(c *fiber.Ctx) error {

	// se debe crear un test en la base de datos.
	// cada test debe tener un total de 10 preguntas.
	// se debe recupear el modulo por el id
	// se debe generar las 10 preguntas para realizar el test.
	// si el usuario no completa las 10 preguntas perdera todo el progreso.

	claims := utils.GetClaims(c)

	idModule, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear el id del modulo",
		})
	}

	testid, err := data.GenerateTestForStudent(claims.UserAPI.ID, uint(idModule))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"testid": testid,
	})
}

func (h *ModuleHandler) GetTest(c *fiber.Ctx) error {

	testId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al parsear el id del test",
		})
	}

	test, err := data.TestByID(uint(testId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(test)

}

func (h *ModuleHandler) ValidationAnswerForTestModule(c *fiber.Ctx) error {
	idquestion, err := c.ParamsInt("answer_user_id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   "Error al recuperar el id de la pregunta",
		})
	}

	var answer types.Answer
	if err := c.BodyParser(&answer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	// Evalular la respuesta del estudiante.
	// Recuperar la answer_user que esta en la base de datos.
	answerUserDB, err := data.GetAnswerUserByID(uint(idquestion))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	// establecemos la nueva respuesta del estudiante.
	answerUserDB.Answer = data.Answer{
		TrueOrFalse:    answer.TrueOrFalse,
		TextOpcions:    answer.TextOpcions,
		TextToComplete: answer.TextToComplete,
	}

	// Evaluación de la pregunta.
	answerUserDB.Responded = true
	switch answerUserDB.Question.TypeQuestion {
	case "true_or_false":
		answerUserDB.IsCorrect = answerUserDB.Question.CorrectAnswer.TrueOrFalse == answerUserDB.Answer.TrueOrFalse
		if answerUserDB.IsCorrect {
			answerUserDB.Score = 10
			answerUserDB.Feedback = "Respuesta correcta"
		} else {
			answerUserDB.Score = 0
			answerUserDB.Feedback = "Respuesta incorrecta"
		}
	case "multi_choice_text":
		// cuanso es multi_chose_text la respuesat viene por TextOpciones.
		if answerUserDB.Question.Options.SelectMode == "single" {
			// si no tiene ninguna opcion selecionada automaticamente es incorrecta
			if len(answerUserDB.Answer.TextOpcions) < 1 {
				answerUserDB.IsCorrect = false
				answerUserDB.Score = 0
				answerUserDB.Feedback = "Respuesta incorrecta"
			} else {
				answerUserDB.IsCorrect = utils.ContainsString(answerUserDB.Question.CorrectAnswer.TextOpcions, answerUserDB.Answer.TextOpcions[0])
				if answerUserDB.IsCorrect {
					answerUserDB.Feedback = "Respuesta correcta"
					answerUserDB.Score = 10
				} else {
					answerUserDB.Feedback = "Respuesta incorrecta"
					answerUserDB.Score = 0
				}
			}
		} else {
			// en caso de ser multiple

			points := 0
			// en caso de ser multiple selección se evalua la respuesta
			for _, correctAnswer := range answerUserDB.Question.CorrectAnswer.TextOpcions {
				if utils.ContainsString(answerUserDB.Answer.TextOpcions, correctAnswer) {
					points++
				}
			}

			answerUserDB.IsCorrect = points == len(answerUserDB.Question.CorrectAnswer.TextOpcions)
			// se calcula el puntaje
			pointsForEachCorrectAnswer := 10 / float32(len(answerUserDB.Question.CorrectAnswer.TextOpcions))
			answerUserDB.Score = pointsForEachCorrectAnswer * float32(points)

			if points == 0 {
				answerUserDB.Feedback = "Respuesta incorrecta"
			} else if points < len(answerUserDB.Question.CorrectAnswer.TextOpcions) {
				count := len(answerUserDB.Question.CorrectAnswer.TextOpcions) - points
				answerUserDB.Feedback = fmt.Sprintf("Te faltó seleccionar %d", count)
			} else {
				answerUserDB.Feedback = "Respuesta correcta"
			}
			//if answerUserDB.IsCorrect {
			//	answerUserDB.Feedback = "Respuesta correcta"
			//} else {
			//	answerUserDB.Feedback = "Respuesta incorrecta"
			//}
		}

	case "complete_word":

		textToCompleteCorrect := []string(answerUserDB.Question.CorrectAnswer.TextToComplete)

		points := 0

		for _, correctAnswer := range textToCompleteCorrect {
			if utils.ContainsString(answerUserDB.Answer.TextToComplete, correctAnswer) {
				points++
				break
			}
		}

		// Calculamos el puntaje para complete word
		answerUserDB.IsCorrect = points == len(textToCompleteCorrect)
		pointsForEachCorrectAnswer := 10 / len(textToCompleteCorrect)
		answerUserDB.Score = float32(pointsForEachCorrectAnswer * points)
		answerUserDB.Feedback = "Respuesta correcta"

	case "order_word":
		// analizamos que la respuesta del usuario sea igual que las opciones correctas
		textToCompleteCorrect := []string(answerUserDB.Question.CorrectAnswer.TextOpcions)
		textToCompleteUser := []string(answerUserDB.Answer.TextOpcions)

		// si el orden esta correcto automaticamente es correcta, en caso de que
		// una no sea correcta automaticamente es incorrecta.
		for i, correctAnswer := range textToCompleteCorrect {
			if i >= len(textToCompleteUser) {
				answerUserDB.IsCorrect = false
				answerUserDB.Score = 0
				answerUserDB.Feedback = "Respuesta incorrecta"
				break
			}
			if correctAnswer == textToCompleteUser[i] {
				answerUserDB.IsCorrect = true
				answerUserDB.Score = 10
				answerUserDB.Feedback = "Respuesta correcta"
				continue
			} else {
				answerUserDB.IsCorrect = false
				answerUserDB.Score = 0
				answerUserDB.Feedback = "Respuesta incorrecta"
				break
			}
		}

	default:
		answerUserDB.IsCorrect = false
		answerUserDB.Score = 0
		answerUserDB.Feedback = ""
	}

	// Actualizar cambios en la base de datos.
	err = answerUserDB.UpdateAnswerUser()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"score":      answerUserDB.Score,
		"is_correct": answerUserDB.IsCorrect,
		"feedback":   answerUserDB.Feedback,
	})
}

func (h *ModuleHandler) FinishTest(c *fiber.Ctx) error {
	// claims := utils.GetClaims(c)
	testid, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "testid not found",
			"error":   err.Error(),
		})
	}

	// Finalizar el test en la base de datos.
	finishTest, err := data.FinishTest(uint(testid))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "error",
			"error":   err.Error(),
		})
	}

	return c.JSON(finishTest)
}

// GetMyTest recupera todos los test de un usuario en un modulo especifico.
func (h *ModuleHandler) GetMyTest(c *fiber.Ctx) error {
	claims := utils.GetClaims(c)
	idModule, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	tests, err := data.GetMyTest(claims.UserAPI.ID, uint(idModule))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	testsAPI := data.TestsModuleToAPI(tests)

	return c.JSON(testsAPI)
}
