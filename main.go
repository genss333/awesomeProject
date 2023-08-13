package main

import (
	"awesomeProject/controller"
	"awesomeProject/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	// Fiber instance
	app := fiber.New()
	app.Use(service.LogRequests)

	// Auth controller routes
	controller.AuthInitializeRoutes(app)

	// Middleware to protect routes with JWT
	app.Use(service.JWTMiddleware)

	// User controller routes
	controller.UserInitializeRoutes(app)

	fmt.Println("Server listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
