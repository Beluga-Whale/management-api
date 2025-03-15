package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(c*fiber.Ctx) error {
	// NOTE - Get cookies 
	tokenString :=  c.Cookies("jwt")    
	
	// NOTE - Check token it empty
	if tokenString == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":"Unauthorized",
		})
	}
	
	// NOTE - นำ token มาเช็คว่าเป็นอันเดียวกันไหม
	token, err := jwt.ParseWithClaims(tokenString,&JWTClaims{}, func(token *jwt.Token)  (interface{},error){
		return []byte(os.Getenv("JWT_SECRET")),nil
	})

	// NOTE - Check token valid
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message":"Invalid or expired token"})
	}

	claims, ok :=token.Claims.(*JWTClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message":"Invalid token claims"})
	}

	fmt.Println(claims)
	return c.Next()

}