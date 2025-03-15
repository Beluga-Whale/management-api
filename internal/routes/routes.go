package routes

import (
	"github.com/Beluga-Whale/management-api/internal/handlers.go"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler ){
	api := app.Group("/api")
	api.Post("/user/register", userHandler.RegisterUser)
	api.Post("/user/login", userHandler.Login)
	api.Post("/user/logout", userHandler.Logout)

}