package user

import (
	"awesomeProject/business/auth"
	"awesomeProject/business/user"
	"awesomeProject/middleware"
	"awesomeProject/service"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {
	app.Get("/api/user/:offset/:limit", middleware.AuthorizationMiddleware([]string{"user", "admin"}), getUsersHandler)
	app.Get("/api/user/:id", middleware.AuthorizationMiddleware([]string{"user", "admin"}), getUserByIdHandler)
	app.Post("/api/user/", middleware.AuthorizationMiddleware([]string{"user", "admin"}), createUserHandler)
	app.Patch("/api/user/update", middleware.AuthorizationMiddleware([]string{"user", "admin"}), updateUserHandler)
	app.Delete("/api/user/delete/:id", middleware.AuthorizationMiddleware([]string{"user", "admin"}), deleteUserHandler)
	app.Post("/api/user/logout", middleware.AuthorizationMiddleware([]string{"user", "admin"}), logoutHandler)
	app.Post("/api/refreshToken", middleware.AuthorizationMiddleware([]string{"user", "admin"}), RefreshTokenHandler)
}

func getUsersHandler(c *fiber.Ctx) error {
	return user.GetUsers(c)
}

func getUserByIdHandler(c *fiber.Ctx) error {
	return user.GetUserById(c)
}

func createUserHandler(c *fiber.Ctx) error {
	return user.CreateUser(c)
}

func updateUserHandler(c *fiber.Ctx) error {
	return user.UpdateUser(c)
}

func deleteUserHandler(c *fiber.Ctx) error {
	return user.DeleteUser(c)
}

func logoutHandler(c *fiber.Ctx) error {
	return auth.Logout(c)
}

func RefreshTokenHandler(c *fiber.Ctx) error {
	return service.RefreshToken(c)
}
