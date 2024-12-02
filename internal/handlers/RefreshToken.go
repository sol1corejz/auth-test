package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/auth/internal/models"
	"github.com/sol1corejz/auth/internal/services"
	"github.com/sol1corejz/auth/internal/storage"
	"log"
	"time"
)

func sendEmailNotification(email, userIP string) error {
	log.Printf("Sending email to %s: IP address changed to %s", email, userIP)
	return nil
}

func RefreshHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case <-ctx.Done():
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "Request timed out",
		})
	default:
		userIP := c.IP()
		accessToken := c.Cookies("accessToken")
		refreshToken := c.Cookies("refreshToken")

		userID, err := services.GetUserID(accessToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"token": accessToken,
			})
		}

		accessData, err := storage.GetUserAccessData(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot get user access data",
			})
		}

		isRefreshValid := services.IsRefreshValid(refreshToken, accessData.RefreshToken)
		if !isRefreshValid {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "invalid refresh token",
			})
		}

		if accessData.UserIP != userIP {
			err = sendEmailNotification("ilyaparunov529@gmail.com", userIP)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to send email notification",
				})
			}
		}

		authTokens, err := services.GenerateTokens(userID, userIP)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		data := models.UserAccessData{
			UserID:       userID,
			UserIP:       userIP,
			RefreshToken: authTokens.HashedRefreshToken,
		}
		err = storage.UpdateUserAccessData(ctx, data)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": "tokens updated",
		})
	}
}
