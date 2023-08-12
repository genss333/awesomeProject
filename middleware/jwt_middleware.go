package middleware

import (
	"awesomeProject/service"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

var secretKey = service.JWTService{
	SecretKey: []byte("8Zz5tw0Ion3XPZZfN0NOml3z9FMultiwordR9fp6ryDIoGRM8STEPHA6iHsc0fb"),
}

func JWTMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing authorization T",
			"status":  fiber.StatusUnauthorized,
			"path":    c.Path(),
		})
	}

	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token format",
			"status":  fiber.StatusUnauthorized,
			"path":    c.Path(),
		})
	}

	jwtConfig := jwtware.Config{
		SigningKey: secretKey.SecretKey,
	}
	jwtMiddleware := jwtware.New(jwtConfig)

	return jwtMiddleware(c)
}

func LogRequests(c *fiber.Ctx) error {
	log.Println(c.Method(), c.Path())
	return c.Next()
}
