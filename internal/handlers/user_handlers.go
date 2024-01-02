package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/resendlabs/resend-go"
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

	// generá el JWT para el usuario.
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

func (h *UserHandler) HandlerResetPassword(c *fiber.Ctx) error {
	// requiere el correo electronico
	var resetPassword types.ResetPassword
	if err := c.BodyParser(&resetPassword); err != nil {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa los datos de la petición",
		})
	}

	ok, user := data.ExisteEmail(resetPassword.Email)

	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "El correo no existe en la base de datos",
		})
	}

	// Generar el jwt para el usuario
	claims := types.UserClaims{
		UserAPI: user,
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

	// guardamos en la db estos datos
	err = data.SaveResetPassword(user.ID, resetPassword.Email, ss)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error al guardar el token", "data": err.Error()})
	}

	apikey := h.config.GetString("KEY_RESEND")
	client := resend.NewClient(apikey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{resetPassword.Email},
		Subject: "Reestablece tu contraseña",
		Html: fmt.Sprintf(`
		<h1>Reestablece tu contraseña</h1>
		<p>Haga click en el siguiente enlace para reestablecer su contraseña</p>
		<p>Token: %s</p>
		<p>El token tiene una duracion de 24 horas.</p>
	`, ss),
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al enviar el correo electronico",
		})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Revisa tu correo electronico"})
}

func (h *UserHandler) HandlerChangePassword(c *fiber.Ctx) error {
	var changePassword types.ChangePassword
	if err := c.BodyParser(&changePassword); err != nil {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa los datos de la petición",
		})
	}

	// Parsear los datos del token
	secret := h.config.GetString("JWT_SECRET")

	token, err := jwt.ParseWithClaims(changePassword.Token, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa los datos de la petición",
		})
	}
	claims, ok := token.Claims.(*types.UserClaims)
	if !ok || !token.Valid {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa los datos de la petición",
		})
	}

	// Revisar si el token no ha sido utilizado.
	used, _, err := data.TokenIsUsed(changePassword.Token)
	if err != nil {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa los datos de la petición",
		})
	}

	if used {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "El token ya fue utilizado",
		})
	}

	// Actualizar la contraseña.
	err = data.UpdatePassword(claims.UserAPI.ID, changePassword.Password)
	if err != nil {
		log.Println("Error al resetear la constraseña", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "No se pudo actuailzar la contraseña",
		})
	}
	// Establecer que el token ya se utilizo.
	err = data.SetTokenUsed(changePassword.Token)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Contraseña actualizada"})
}
