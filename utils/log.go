package utils

import (
	"github.com/gofiber/fiber/v2"
)

func RespondJson(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"path":    c.Path(),
		"message": message,
		"status":  statusCode,
	})
}
