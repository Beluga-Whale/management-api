package middleware

import (
	"fmt"

	"github.com/Beluga-Whale/management-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(c*fiber.Ctx) error {
	jwtUtil := utils.NewJwt()

	// NOTE - Get cookies 
	tokenString :=  c.Cookies("jwt")    
	
	// NOTE - Check token it empty
	if tokenString == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":"Unauthorized",
		})
	}
	
	claims,err := jwtUtil.ParseJWT(tokenString)

	c.Locals("userEmail", claims)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":"Invalid token claims",
		})
	}

	fmt.Println(claims)
	return c.Next()

}