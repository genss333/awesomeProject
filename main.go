package main

import (
	"awesomeProject/controller/auth"
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
	database.MigrationTable()

	// Auth controller routes
	auth.InitializeRoutes(app)

	// Security configuration
	security.ApplySecurityConfiguration(app)

	fmt.Println("Server listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
