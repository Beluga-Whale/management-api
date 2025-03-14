package services

import (
	"errors"
	"fmt"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(user *models.Users) error {
	if user.Email == ""{
		return errors.New("Email is required")
	}
	// NOTE - Check email
	result, err := s.userRepo.FindByEmail(user.Email)

	if err != nil{
		return fmt.Errorf("Fail To Check Email : %w",err)
	}


	if result != nil {
		return errors.New("Email has already been used")
	}
	// NOTE - Check password length
	if len(user.Password)<6 {
		return errors.New("Password ")
	}

	return s.userRepo.CreateUser(user)
}

// func (s *UserService) Login(user *models.Users) error {
// 	if user.Email =="" || user.Password =="" {
// 		return errors.New("Email or Password is valid")
// 	} 


// }

