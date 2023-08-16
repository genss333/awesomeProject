package controller

import (
	"awesomeProject/business"
	"awesomeProject/middleware"
	"awesomeProject/service"
	"github.com/gofiber/fiber/v2"
)

func UserInitializeRoutes(app *fiber.App) {
	app.Get("/api/user/:offset/:limit", middleware.AuthorizationMiddleware([]string{"user", "admin"}), getUsersHandler)
	app.Get("/api/user/:id", middleware.AuthorizationMiddleware([]string{"user", "admin"}), getUserByIdHandler)
	app.Post("/api/user/", middleware.AuthorizationMiddleware([]string{"user", "admin"}), createUserHandler)
	app.Patch("/api/user/update", middleware.AuthorizationMiddleware([]string{"user", "admin"}), updateUserHandler)
	app.Delete("/api/user/delete/:id", middleware.AuthorizationMiddleware([]string{"user", "admin"}), deleteUserHandler)
	app.Post("/api/user/logout", middleware.AuthorizationMiddleware([]string{"user", "admin"}), logoutHandler)
	app.Post("/api/refreshToken", middleware.AuthorizationMiddleware([]string{"user", "admin"}), RefreshTokenHandler)
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

func deleteUserHandler(c *fiber.Ctx) error {
	return business.DeleteUser(c)
}

func logoutHandler(c *fiber.Ctx) error {
	return business.Logout(c)
}

func RefreshTokenHandler(c *fiber.Ctx) error {
	return service.RefreshToken(c)
}
