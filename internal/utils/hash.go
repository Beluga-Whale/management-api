package utils

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type HashInterface interface {
	CheckPassword(user *models.Users, password string) bool
}

type Hash struct{}

func NewHash() *Hash{
	return &Hash{}
}

func (h *Hash) CheckPassword(user *models.Users, password string) bool {
	if user == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}