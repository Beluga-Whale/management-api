package routes

import (
	"github.com/Beluga-Whale/management-api/internal/handlers.go"
	"github.com/Beluga-Whale/management-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler ){
	api := app.Group("/api")
	api.Post("/user/register", userHandler.RegisterUser)
	api.Post("/user/login", userHandler.Login)
	api.Post("/user/logout", userHandler.Logout)

	// NOTE - Protect routes by authMiddleware
	api.Use(middleware.AuthMiddleware)
	api.Get("/task",func (c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message":"TEST"})
	})

}