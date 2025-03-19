package handlers

import (
	"time"

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

	// NOTE - Call service Register
	err := h.userService.RegisterUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message":"User registered success"})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	user := new(models.Users)
	if err:= c.BodyParser(user); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"In valid request"})
	}

	// NOTE - Call Service login
	token, userDetail ,err := h.userService.Login(user)

	if err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	// NOTE - Set cookie
	c.Cookie(&fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*72),
		HTTPOnly: true,
		Secure:false,
		SameSite: fiber.CookieSameSiteNoneMode, 
		
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Login success",
		"token":token,
		"user":fiber.Map{
			"id": userDetail.ID,
			"email":userDetail.Email,
			"name": userDetail.Name,
		},
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie();

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Logout Success",
	})
}