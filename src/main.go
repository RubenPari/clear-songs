package main

import (
	"os"

	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	errEnv := godotenv.Load("../.env")
	if errEnv != nil {
		panic(errEnv)
	}

	routes.SetUpRoutes(app)

	errServer := app.Listen(":" + os.Getenv("PORT"))
	if errServer != nil {
		panic(errServer)
	}
}
