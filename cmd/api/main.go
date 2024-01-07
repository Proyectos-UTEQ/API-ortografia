package main

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/internal/handlers"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	config := viper.New()

	// Read environment variables
	config.AutomaticEnv()

	config.SetDefault("APP_PORT", "3000")
	config.SetDefault("APP_ENV", "development")

	// Read the config file
	config.SetConfigName("config")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	// config.AddConfigPath("/workspaces/api-ortografia")

	// Load the config
	err := config.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	// Connect to the database
	database := db.ConnectDB(config)

	// Migrate the schema
	err = database.AutoMigrate(
		&data.User{},
		&data.ResetPassword{},
		&data.Module{},
	)
	if err != nil {
		fmt.Println(err)
	}

	// Create fiber app
	app := fiber.New()

	// Create handlers
	userHandler := handlers.NewUserHandler(config)
	jwtHandler := handlers.NewJWTHandler(config)
	moduleHandler := handlers.NewModuleHandler(config)

	api := app.Group("/api")

	auth := api.Group("/auth")
	// Routes for auth users
	auth.Post("/sign-in", userHandler.HandlerSignin)
	auth.Post("/sign-up", userHandler.HandlerSignup)

	// se encarga de enviar el correo electronico al usuario
	auth.Post("/reset-password", userHandler.HandlerResetPassword)

	// se encarga de actulizar la constraseña del usuario
	// esto debe resivir un token.
	auth.Put("/change-password", userHandler.HandlerChangePassword)

	// Ejemplo de rutas protegidas.
	api.Get("/protegida", jwtHandler.JWTMiddleware, handlers.Authorization("admin", "teacher"), func(c *fiber.Ctx) error {
		claims := utils.GetClaims(c)
		fmt.Println(claims.UserAPI)
		return c.SendString("ruta protegida, has tenido acceso " + claims.UserAPI.FirstName)
	})

	module := api.Group("/module", jwtHandler.JWTMiddleware) // solo con JWT se tiene acceso.

	// Routes for modules
	module.Post(
		"/",
		handlers.Authorization("teacher", "admin"),
		moduleHandler.CreateModuleForTeacher)

	// Lista todos los modulos.
	module.Get(
		"/teacher",
		handlers.Authorization("teacher", "admin"),
		moduleHandler.GetModulesForTeacher,
	)

	module.Get(
		"/",
		moduleHandler.GetModules,
	)

	upload := api.Group("/upload")
	uploadHandler := handlers.NewUploadHandler(config)

	upload.Post("/", jwtHandler.JWTMiddleware, uploadHandler.UploadFiles)
	upload.Static("/", "./uploads")

	app.Listen(":" + config.GetString("APP_PORT"))
}
