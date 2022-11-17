package main

import (
	"os"

	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	authMO.LoadEnv(3)

	routes.SetUpRoutes(app)

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}
}
