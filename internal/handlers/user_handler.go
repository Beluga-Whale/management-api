package handlers

import (
	"time"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler{
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
		Secure:true,
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

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userEmail := c.Locals("userEmail").(string)

	user, err := h.userService.GetUserByEmail(userEmail)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":user,
	})
}

func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	user := new(models.Users)

	err := c.BodyParser(user)
		// NOTE - get ID From Params
		idStr:= c.Params("id")
		if idStr == ""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Task ID is required",
			})
		}
	// NOTE - ดึง Email จาก cookie
	emailCookie := c.Cookies("jwt")
	if emailCookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request",
		})
	}

	if err :=h.userService.UpdateUserById(idStr, emailCookie, user); err !=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error :": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":user,
	})
}

