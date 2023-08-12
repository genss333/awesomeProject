package controller

import (
	"awesomeProject/business"
	fiber "github.com/gofiber/fiber/v2"
)

func AuthInitializeRoutes(app *fiber.App) {
	app.Post("/api/login", loginHandler)
}

func loginHandler(c *fiber.Ctx) error {
	return business.Login(c)
}
