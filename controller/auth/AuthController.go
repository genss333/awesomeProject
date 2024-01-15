package auth

import (
	"awesomeProject/business/auth"
	fiber "github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {
	app.Post("/api/login", loginHandler)
	app.Post("/api/register", RegisterHandler)
}

func loginHandler(c *fiber.Ctx) error {
	return auth.Login(c)
}

func RegisterHandler(c *fiber.Ctx) error {
	return auth.Register(c)
}
