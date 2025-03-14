package handlers

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler{
	return &UserHandler{userService:userService}
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	user := new(models.Users)
	if err:= c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid request"})
	}

	err := h.userService.RegisterUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message":"User registered success"})
}