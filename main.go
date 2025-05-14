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
	"github.com/Beluga-Whale/management-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectDB()

	// NOTE - Fiber
	app := fiber.New()

	// NOTE - Use cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:3001, http://13.54.74.23:3000",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// NOTE - Create Repository
	userRepo := repositories.NewUserRepository(config.DB)
	taskRepo := repositories.NewTaskRepository(config.DB)

	hashUtil := utils.NewHash()
	jwtUtil := utils.NewJwt()
	// NOTE - Create Service
	userService := services.NewUserService(userRepo,hashUtil,jwtUtil)
	taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

	// NOTE - Handler
	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// NOTE - Route 
	routes.SetupRoutes(app,userHandler,taskHandler)


	port := os.Getenv("PORT_API")

	if port =="" {
		port =":8080"
	}
	 addr := fmt.Sprintf(":%s", port)

    fmt.Printf("Server running on port %s\n", port)
    // NOTE -เช็ค error จาก Listen
    if err := app.Listen(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }

}