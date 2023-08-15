package controller

import (
	"awesomeProject/business"
	"awesomeProject/service"
	fiber "github.com/gofiber/fiber/v2"
)

func AuthInitializeRoutes(app *fiber.App) {
	app.Post("/api/login", loginHandler)
	app.Post("/api/register", RegisterHandler)
	app.Post("/api/refresh", RefreshTokenHandler)
}

func loginHandler(c *fiber.Ctx) error {
	return business.Login(c)
}

func RegisterHandler(c *fiber.Ctx) error {
	return business.CreateUser(c)
}

func RefreshTokenHandler(c *fiber.Ctx) error {
	return service.RefreshToken(c)
}
