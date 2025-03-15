package services

import (
	"errors"
	"fmt"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/utils"
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
		return errors.New("Password must more 6 char ")
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) Login(user *models.Users) (string,*models.Users,error) {
	if user.Email =="" || user.Password =="" {
		return "",nil,errors.New("Email or Password is required")
	} 
		
	// NOTE - find User by Email
	dbUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return "",nil,errors.New("Email not found")
	}

	if !utils.CheckPassword(dbUser ,user.Password) {
		return "",nil,errors.New("Invalid Email or Password")
	}

	token, err := utils.GenerateJWT(dbUser.Email)
	
	if err != nil{
		return "",nil,fmt.Errorf("Failed to generate token: %w", err)
	}

	return token,dbUser,nil

}

