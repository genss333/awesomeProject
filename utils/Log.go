package utils

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func RespondJson(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"path":    c.Path(),
		"message": message,
		"status":  statusCode,
	})
}

func LogRequests(c *fiber.Ctx) error {
	log.Println(c.Method(), c.Path())
	return c.Next()
}
