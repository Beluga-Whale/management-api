package routes

import (
	"github.com/Beluga-Whale/management-api/internal/handlers"
	"github.com/Beluga-Whale/management-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler ){
	api := app.Group("/api")
	api.Post("/user/register", userHandler.RegisterUser)
	api.Post("/user/login", userHandler.Login)
	api.Post("/user/logout", userHandler.Logout)

	// NOTE - Protect routes by authMiddleware
	api.Use(middleware.AuthMiddleware)

	// NOTE - Task routes
	api.Post("/task", taskHandler.CreateTask)
	api.Get("/task", taskHandler.GetAllTask)
	api.Get("/task/complete", taskHandler.GetCompleteTask)
	api.Get("/task/pending", taskHandler.GetPendingTask)
	api.Get("/task/overdue", taskHandler.GetOverdueTask)
	
	api.Get("/task/:id", taskHandler.FindTaskById)
	api.Put("/task/:id",taskHandler.UpdateTask)
	api.Delete("task/:id", taskHandler.DeleteTask)

	// NOTE - User routes
	api.Get("/user", userHandler.GetUser)

}