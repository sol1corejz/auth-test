package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/auth/internal/handlers"
)

func main() {

	app := fiber.New()

	app.Get("/get-tokens/:id", handlers.GetTokensHandler)
	app.Post("/fefresh-token", handlers.RefreshHandler)

	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
