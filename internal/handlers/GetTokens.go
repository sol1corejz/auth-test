package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/auth/internal/services"
)

func GetTokensHandler(c *fiber.Ctx) error {
	userUUID := c.Params("id")
	userIP := c.IP()

	fmt.Println(userUUID, userIP)

	authTokens, err := services.GenerateTokens(userUUID, userIP)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"tokens": authTokens,
	})
}
