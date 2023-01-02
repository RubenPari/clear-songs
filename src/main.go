package main

import (
	"os"

	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.SetUpRoutes(app)

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}
}
