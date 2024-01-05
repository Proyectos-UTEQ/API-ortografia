package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTHandler struct {
	config *viper.Viper
}

func NewJWTHandler(config *viper.Viper) *JWTHandler {
	return &JWTHandler{
		config: config,
	}
}

func (h *JWTHandler) JWTMiddleware(c *fiber.Ctx) error {
	// recuperar el token
	auth := c.Get("Authorization")
	authArray := strings.Split(auth, " ")
	if len(authArray) < 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token no valido"})
	}

	tokenString := authArray[1]

	secret := h.config.GetString("APP_JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token no valido"})
	}

	claims, ok := token.Claims.(*types.UserClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token no valido"})
	}

	c.Locals("user", claims)

	return c.Next()
}

// controla la autorización de los usuario por rol.
func Authorization(allowedRoles ...string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		claims := utils.GetClaims(c)
		for _, role := range allowedRoles {
			if role == claims.TypeUser {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "No autorizado"})
	}

}
