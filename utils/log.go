package utils

import (
	"github.com/gofiber/fiber/v2"
)

func RespondWithError(c *fiber.Ctx, statusCode int, message string) {
	err := c.Status(statusCode).JSON(fiber.Map{
		"path":    c.Path(),
		"message": message,
		"status":  statusCode,
	})
	if err != nil {
		return
	}
}

func RespondWithSuccess(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"path":    c.Path(),
		"message": message,
		"status":  statusCode,
	})
}
