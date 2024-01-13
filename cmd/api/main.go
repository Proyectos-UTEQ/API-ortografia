package main

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/internal/handlers"
	"Proyectos-UTEQ/api-ortografia/internal/services"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	config.AddConfigPath("/etc/secrets/")
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
		&data.Subscription{},
	)
	if err != nil {
		fmt.Println(err)
	}

	// Create fiber app
	app := fiber.New()

	// configuració de cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

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

	// Actuliza el modulo
	module.Put(
		"/:id",
		handlers.Authorization("teacher", "admin"),
		moduleHandler.UpdateModule,
	)

	// Lista todos los modulos.
	module.Get(
		"/teacher",
		handlers.Authorization("teacher", "admin"),
		moduleHandler.GetModulesForTeacher,
	)

	// Lista todos los modulos.
	module.Get(
		"/",
		moduleHandler.GetModules,
	)

	// Recupera todos los modulos y ademas indica si el usuario esta suscrito o no.
	module.Get("/with-is-subscribed", moduleHandler.GetModuleWithIsSubscribed)

	module.Post("/subscribe", moduleHandler.Subscribe)
	module.Get("/subscribed", moduleHandler.Subscriptions)

	// Listar todos los estudiantes que estan suscritos a un modulo.
	module.Get("/:id/students", moduleHandler.GetStudents)

	// Routes for modules
	// Crea un modulo.
	module.Post(
		"/",
		handlers.Authorization("teacher", "admin"),
		moduleHandler.CreateModuleForTeacher)

	module.Get("/:id", moduleHandler.GetModuleByID)

	// Routes for upload
	upload := api.Group("/uploads")
	uploadHandler := handlers.NewUploadHandler(config)

	upload.Post("/", jwtHandler.JWTMiddleware, uploadHandler.UploadFiles)
	upload.Static("/", "./uploads")

	go services.TelegramBot(config)
	app.Listen(":" + config.GetString("APP_PORT"))
}
