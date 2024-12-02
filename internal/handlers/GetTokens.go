package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/auth/internal/models"
	"github.com/sol1corejz/auth/internal/services"
	"github.com/sol1corejz/auth/internal/storage"
	"time"
)

func GetTokensHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case <-ctx.Done():
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "Request timed out",
		})
	default:
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "userID is required",
			})
		}
		userIP := c.IP()

		authTokens, err := services.GenerateTokens(userID, userIP)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		data := models.UserAccessData{
			UserID:       userID,
			UserIP:       userIP,
			RefreshToken: authTokens.HashedRefreshToken,
		}

		err = storage.CreateUserAccessData(ctx, data)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "accessToken",
			Value:    authTokens.AccessToken,
			HTTPOnly: true,
		})
		c.Cookie(&fiber.Cookie{
			Name:     "refreshToken",
			Value:    authTokens.RefreshToken,
			HTTPOnly: true,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Tokens pair successfully created",
		})
	}
}
