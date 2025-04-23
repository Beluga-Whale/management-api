package services

import (
	"errors"
	"fmt"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/utils"
)

type UserServiceInterface interface {
	RegisterUser(user *models.Users) error
	Login(user *models.Users) (string,*models.Users,error)
	GetUserByEmail(email string) (*models.Users, error)
	UpdateUserById(idStr string, emailCookie string, updatedUserValue *models.Users) (error)
}

type UserService struct {
	userRepo repositories.UserRepositoryInterface
	hashUtil utils.HashInterface
	jwtUtil utils.JwtInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface,hashUtil utils.HashInterface, jwtUtil utils.JwtInterface) *UserService {
	return &UserService{userRepo: userRepo,hashUtil:hashUtil, jwtUtil: jwtUtil  }
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

	if !s.hashUtil.CheckPassword(dbUser ,user.Password) {
		return "",nil,errors.New("Invalid Email or Password")
	}

	token, err := s.jwtUtil.GenerateJWT(dbUser.Email)
	
	if err != nil{
		return "",nil,fmt.Errorf("Failed to generate token: %w", err)
	}

	return token,dbUser,nil

}

func (s *UserService) GetUserByEmail(email string) (*models.Users, error) {


	return s.userRepo.FindByEmail(email)
}

func (s *UserService) UpdateUserById(idStr string, emailCookie string, updatedUserValue *models.Users) (error) {
	// NOTE - Check idStr
	if idStr == "" {
		return errors.New("Id is required")
	}

	// NOTE - Decode Jwt in cookie เพื่อดึง Eamil
	email,err :=s.jwtUtil.ParseJWT(emailCookie)

	if err != nil{
		return fmt.Errorf("Fail To Check Email : %w",err)
	}

	// NOTE - หา User จาก Email เพื่อเอา UserID 
	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return  errors.New("User not found")
	}


	// NOTE -หา Task By ID
	userID,err:= s.userRepo.FindUserById(idStr)

	if err != nil {
		return  fmt.Errorf("failed to find task by ID: %w", err)
	}

	// NOTE - มาเช็คว่าผู้ใช้เป็นเจ้าของ Task ไหม
	if userID.ID != user.ID {
		return  errors.New("you do not have permission to access this task")
	}
	if	err :=s.userRepo.UpdateUserById(updatedUserValue,userID.ID); err != nil {
		return fmt.Errorf("Error : %w",err)
	}
	
	return nil


}