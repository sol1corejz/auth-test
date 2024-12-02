package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/auth/cmd/config"
	"github.com/sol1corejz/auth/internal/handlers"
	"github.com/sol1corejz/auth/internal/storage"
)

func main() {

	config.ParseFlags()

	if err := storage.InitStorage(); err != nil {
		fmt.Println(err)
		return
	}

	app := fiber.New()

	app.Get("/get-tokens/:id", handlers.GetTokensHandler)
	app.Post("/refresh-token", handlers.RefreshHandler)

	err := app.Listen(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}
