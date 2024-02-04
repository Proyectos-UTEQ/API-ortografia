package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
)

type ClassesHandler struct {
	config *viper.Viper
}

// NewClassesHandler crea un nuevo handler de clases.
func NewClassesHandler(config *viper.Viper) *ClassesHandler {
	return &ClassesHandler{
		config: config,
	}
}

func (h *ClassesHandler) NewClasses(c *fiber.Ctx) error {
	// Recuperamos los claims del usuarios
	claims := utils.GetClaims(c)

	// Obtenemos el ID del usuario que crea la clase.
	idCreatorUser := claims.UserAPI.ID

	// Parseamos el body
	var classAPI types.Class
	if err := c.BodyParser(&classAPI); err != nil {
		log.Println("Error al parsear el body")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Establecemos el creador de la clase.
	classAPI.CreatedByID = idCreatorUser

	// generamos el código de la clase
	classAPI.Code = uuid.NewString()

	// Validamos la clase
	err := classAPI.ValidateNewClass()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Guardamos la clase en la base de datos.

	idClass, err := data.RegisterClass(classAPI)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Recuperamos los datos de una clase por el ID.
	class, err := data.GetClassByID(idClass)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// convertir la clase en un json
	classAPI = data.ClassToAPI(class)

	return c.Status(fiber.StatusOK).JSON(classAPI)
}

func (h *ClassesHandler) GetClassesByTeacher(c *fiber.Ctx) error {
	// Recuperar todas las clases basándonos en el profesor que las solicita.
	//claims := utils.GetClaims(c)

	idTeacher, err := c.ParamsInt("id")

	classes, err := data.GetClassesByTeacherID(uint(idTeacher))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	classesAPI := data.ClassesToAPI(classes)
	return c.Status(fiber.StatusOK).JSON(classesAPI)
}

func (h *ClassesHandler) GetClassesArchivedByTeacher(c *fiber.Ctx) error {
	// Recuperar todas las clases basándonos en el profesor que las solicita.
	//claims := utils.GetClaims(c)

	idTeacher, err := c.ParamsInt("id")

	classes, err := data.GetClassesArchivedByTeacherID(uint(idTeacher))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	classesAPI := data.ClassesToAPI(classes)
	return c.Status(fiber.StatusOK).JSON(classesAPI)
}

// ArchiveClassByID Archivar clase por el ID.
func (h *ClassesHandler) ArchiveClassByID(c *fiber.Ctx) error {

	// Recuperamos el ID de la clase.
	idClass, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Actualizamos el registro de la clase
	err = data.ArchiveClass(uint(idClass))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Respondemos con un OK.
	return c.SendStatus(fiber.StatusOK)
}
