package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Beluga-Whale/management-api/config"
	"github.com/Beluga-Whale/management-api/internal/handlers"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/routes"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// NOTE - Connect DB
	config.ConnectDB()

	// NOTE - Fiber
	app := fiber.New()

	// NOTE - Use cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:3001",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// NOTE - Create Repository
	userRepo := repositories.NewUserRepository(config.DB)
	taskRepo := repositories.NewTaskRepository(config.DB)

	// NOTE - Create Service
	userService := services.NewUserService(userRepo)
	taskService := services.NewTaskService(taskRepo,userRepo)

	// NOTE - Handler
	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// NOTE - Route 
	routes.SetupRoutes(app,userHandler,taskHandler)


	port := os.Getenv("PORT_API")

	if port =="" {
		port =":8080"
	}
	fmt.Println("Server running on port",string(port))

	app.Listen(port)

}