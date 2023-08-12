package main

import (
	"awesomeProject/controller"
	"awesomeProject/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	// Fiber instance
	app := fiber.New()
	app.Use(middleware.LogRequests)

	// Auth controller routes
	controller.AuthInitializeRoutes(app)

	// Middleware to protect routes with JWT
	app.Use(middleware.JWTMiddleware)

	// User controller routes
	controller.UserInitializeRoutes(app)

	fmt.Println("Server listening on :8080")
	log.Fatal(app.Listen(":8080"))
}