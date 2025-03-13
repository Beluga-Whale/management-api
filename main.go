package main

import (
	"github.com/Beluga-Whale/management-api/config"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// NOTE - Connect DB
	config.ConnectDB()

	// NOTE - Fiber
	app := fiber.New()

	// NOTE - Create Repository
	userRepo := repositories.NewUserRepository(config.DB)
	taskRepo := repositories.NewUserRepository(config.DB)

	
}