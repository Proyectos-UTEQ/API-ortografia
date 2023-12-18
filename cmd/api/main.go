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

	config.AutomaticEnv()

	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	database := db.InitDB(config)
	// Migrate the schema
	err = database.AutoMigrate(&data.User{})
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()
	api := app.Group("/api")

	userHandler := handlers.NewUserHandler()

	api.Get("/users", userHandler.Login)

	app.Listen(":3000")
}
