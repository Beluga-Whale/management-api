package utils

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(user *models.Users, password string) bool {
	if user == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}