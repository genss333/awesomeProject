package controller

import (
	"awesomeProject/business"
	fiber "github.com/gofiber/fiber/v2"
)

func UserInitializeRoutes(app *fiber.App) {
	app.Get("/api/users", getUsersHandler)
	app.Get("/api/users/:id", getUserByIdHandler)
}

func getUsersHandler(c *fiber.Ctx) error {
	return business.GetUsers(c)
}

func getUserByIdHandler(c *fiber.Ctx) error {
	return business.GetUserById(c)
}
