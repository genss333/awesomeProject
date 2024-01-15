package security

import (
	"awesomeProject/controller/user"
	"awesomeProject/middleware"
	"github.com/gofiber/fiber/v2"
)

func ApplySecurityConfiguration(app *fiber.App) {
	// Middleware to protect routes with JWT
	app.Use(middleware.AuthenticationMiddleware)

	// User controller routes
	user.InitializeRoutes(app)
}
