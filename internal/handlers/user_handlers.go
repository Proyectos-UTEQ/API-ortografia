package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type UserHandler struct {
	config *viper.Viper
}

func NewUserHandler(config *viper.Viper) *UserHandler {
	return &UserHandler{
		config: config,
	}
}

// HandlerSignin inici de sesion para los usuarios.
func (h *UserHandler) HandlerSignin(c *fiber.Ctx) error {
	// parseamso los datos.
	var login types.Login
	if err := c.BodyParser(&login); err != nil {
		log.Println("Error al iniciar sesion", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa tu solicitud",
		})
	}

	// realizamos la autenticacion del usuario
	user, ok, err := data.Login(login)

	if !ok || err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error al iniciar sesion", "data": err.Error()})
	}

	// TODO: generar el JWT para el usuario.
	claims := types.UserClaims{
		UserAPI: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	secret := h.config.GetString("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error al generar el token", "data": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": ss,
		"user":  user,
	})
}

// HandlerSignup crea un nuevo usuario.
func (h *UserHandler) HandlerSignup(c *fiber.Ctx) error {

	// parseamos los datos
	var user types.UserAPI
	if err := c.BodyParser(&user); err != nil {
		log.Println("Error al registrar usuario", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa tu solicitud",
		})
	}

	// Validar datos para registro inicial.
	resp, err := types.Validate(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error en la validacion de datos",
			"data":    resp,
		})
	}

	// Crea el usuario en la base de datos.
	err = data.Register(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al registrar usuario",
			"data":    err,
		})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User created", "data": user})
}
