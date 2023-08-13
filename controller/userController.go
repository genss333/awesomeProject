package controller

import (
	"awesomeProject/business"
	"github.com/gofiber/fiber/v2"
)

func UserInitializeRoutes(app *fiber.App) {
	app.Get("/api/user/", getUsersHandler)
	app.Get("/api/user/:id", getUserByIdHandler)
	app.Post("/api/user/", createUserHandler)
	app.Patch("/api/user/update", updateUserHandler)
	app.Post("/api/user/logout", logoutHandler)
}

func getUsersHandler(c *fiber.Ctx) error {
	return business.GetUsers(c)
}

func getUserByIdHandler(c *fiber.Ctx) error {
	return business.GetUserById(c)
}

func createUserHandler(c *fiber.Ctx) error {
	return business.CreateUser(c)
}

func updateUserHandler(c *fiber.Ctx) error {
	return business.UpdateUser(c)
}

func logoutHandler(c *fiber.Ctx) error {
	return business.Logout(c)
}
