package main

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/internal/handlers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	config := viper.New()

	// Read environment variables
	config.AutomaticEnv()

	// Read the config file
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	// Load the config
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Connect to the database
	database := db.ConnectDB(config)

	// Migrate the schema
	err = database.AutoMigrate(&data.User{})
	if err != nil {
		fmt.Println(err)
	}

	// Create fiber app
	app := fiber.New()
	api := app.Group("/api")

	userHandler := handlers.NewUserHandler()

	api.Get("/users", userHandler.Login)

	app.Listen(":3000")
}
