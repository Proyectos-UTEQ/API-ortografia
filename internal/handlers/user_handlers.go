package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/services"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"fmt"
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

	if user.URLAvatar == "" {
		user.URLAvatar = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=5952A2&color=fff&size=128", user.FirstName)
	} else {
		user.URLAvatar = h.config.GetString("APP_HOST") + user.URLAvatar
	}

	// generá el JWT para el usuario.
	claims := types.UserClaims{
		UserAPI: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	secret := h.config.GetString("APP_JWT_SECRET")
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
	var userAPI types.UserAPI
	if err := c.BodyParser(&userAPI); err != nil {
		log.Println("Error al registrar usuario", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Revisa tu solicitud",
		})
	}

	// Validar datos para registro inicial.
	resp, err := types.Validate(&userAPI)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error en la validacion de datos",
			"data":    resp,
		})
	}

	// Crea el usuario en la base de datos.
	user, err := data.Register(&userAPI)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al registrar usuario",
			"data":    err.Error(),
		})
	}

	// convertir los datos a un usuario api
	result := data.UserToAPI(*user)

	return c.JSON(fiber.Map{"status": "success", "message": "User created", "data": result})
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

	secret := h.config.GetString("APP_JWT_SECRET")
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

	// make message to send
	messageToSend := fmt.Sprintf(`Hola, %s %s. Haga click en el siguiente enlace para reestablecer su contraseña: http://localhost:3000/reset-password/%s`, user.FirstName, user.LastName, ss)

	emailNotifier := services.NewEmailNotifier(h.config, []string{user.Email}, "Reestablece tu contraseña")
	telegramNotifier := services.NewTelegramNotifier(h.config, user.TelegramID)

	err = utils.ResetPassword(emailNotifier, messageToSend, "https://app-poliword.onrender.com/auth/forgot-password/"+ss)
	if err != nil {
		log.Println(err)
	}
	err = utils.ResetPassword(telegramNotifier, "Presiona el siguiente boton para resetear tu contraseña", "https://app-poliword.onrender.com/auth/forgot-password?token="+ss)
	if err != nil {
		log.Println(err)
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
	secret := h.config.GetString("APP_JWT_SECRET")

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

// HandlerGetUser recupera los datos del usuario en base al token.
func (h *UserHandler) HandlerGetUser(c *fiber.Ctx) error {

	claims := utils.GetClaims(c)

	user, err := data.GetUserByID(claims.UserAPI.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(user)
}

// HandlerUpdateUser actualiza los datos del usuario en base al token.
func (h *UserHandler) HandlerUpdateUser(c *fiber.Ctx) error {
	claims := utils.GetClaims(c)

	var user types.UserAPI
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := user.ValidateUpdateUser()
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Actualizamos en la base de datos.
	err = data.UpdateUserByID(user, claims.UserAPI.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetAllUsers obtiene todos los usuarios
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {

	// Obtenemos todos los usuarios de la DB.
	users, err := data.GetAllUsers()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(users)
}

func (h *UserHandler) ActiveUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = data.ChangeStatus(uint(userID), true)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) BlockedUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = data.ChangeStatus(uint(userID), false)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
