package main

import (
	"awesomeProject/controller"
	"awesomeProject/database"
	"awesomeProject/middleware"
	"awesomeProject/security"
	"awesomeProject/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	// Fiber instance
	app := fiber.New()
	app.Use(utils.LogRequests)
	app.Use(middleware.CorsMiddleware())

	//Migrate database
	database.CreateTables()

	// Auth controller routes
	controller.AuthInitializeRoutes(app)

	// Security configuration
	security.ApplySecurityConfiguration(app)

	fmt.Println("Server listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
