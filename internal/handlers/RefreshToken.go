package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func RefreshHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
