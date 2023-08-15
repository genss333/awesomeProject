package security

import (
	"awesomeProject/controller"
	"awesomeProject/middleware"
	"github.com/gofiber/fiber/v2"
)

func ApplySecurityConfiguration(app *fiber.App) {
	// Middleware to protect routes with JWT
	app.Use(middleware.JWTMiddleware)

	// User controller routes
	controller.UserInitializeRoutes(app)
}
