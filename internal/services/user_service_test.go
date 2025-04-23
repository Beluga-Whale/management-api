package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/Beluga-Whale/management-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	t.Run("Register success",func(t *testing.T) {
		user := &models.Users{
			Email: "Test@gmail.com",
			Password: "Testasdfsdf",
			Name: "tester",
		}
		userRepo := repositories.NewUserRepositoryMock()

		// NOTE - ส่ง email ไปเช็คว่าซ้ำกันในระบบไหม
		userRepo.On("FindByEmail",user.Email).Return(nil,nil)

		// NOTE - สมัครต่อได้
		userRepo.On("CreateUser",mock.Anything).Return(nil)
		
		checkPassword := utils.NewHashMock()
		userService := services.NewUserService(userRepo,checkPassword)

		err :=userService.RegisterUser(user)

		assert.NoError(t, err)

		// NOTE - เช็คว่ามีการ call function ที่เราเรียกจริงไหม
		userRepo.AssertExpectations(t)
	})
	t.Run("Required Email",func(t *testing.T) {
		user := &models.Users{
			Email: "",
			Password: "Testasdfsdf",
			Name: "tester",
		}
		userRepo := repositories.NewUserRepositoryMock()
		checkPassword := utils.NewHashMock()
		userService := services.NewUserService(userRepo,checkPassword)
	
		err :=userService.RegisterUser(user)
		assert.EqualError(t,err,"Email is required")
	})

	t.Run("Error to check email",func(t *testing.T) {
		user := &models.Users{
			Email: "error@mail.com",
			Password: "Testasdfsdf",
			Name: "tester",
		}
		userRepo := repositories.NewUserRepositoryMock()

		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("Fail To Check Email"))

		checkPassword := utils.NewHashMock()
		userService := services.NewUserService(userRepo,checkPassword)

		err := userService.RegisterUser(user)

		assert.EqualError(t,err,"Fail To Check Email : Fail To Check Email")
	})

	t.Run("Email is already",func(t *testing.T) {
		user := &models.Users{
			Email: "already@mail.com",
			Password: "Testasdfsdf",
			Name: "tester",
		}
		userRepo := repositories.NewUserRepositoryMock()

		userRepo.On("FindByEmail",user.Email).Return(user,nil)

		checkPassword := utils.NewHashMock()
		userService := services.NewUserService(userRepo,checkPassword)

		err := userService.RegisterUser(user)

		assert.EqualError(t,err,"Email has already been used")
		// NOTE - เช็คว่ามีการ call function ที่เราเรียกจริงไหม
		userRepo.AssertExpectations(t)
	})

	t.Run("Password less more than 6",func(t *testing.T) {
		user := &models.Users{
			Email: "already@mail.com",
			Password: "Tef",
			Name: "tester",
		}
		userRepo := repositories.NewUserRepositoryMock()

		// NOTE - ส่ง email ไปเช็คว่าซ้ำกันในระบบไหม
		userRepo.On("FindByEmail",user.Email).Return(nil,nil)

		checkPassword := utils.NewHashMock()
		userService := services.NewUserService(userRepo,checkPassword)

		err :=userService.RegisterUser(user)

		assert.EqualError(t, err,"Password must more 6 char ")

		// NOTE - เช็คว่ามีการ call function ที่เราเรียกจริงไหม
		userRepo.AssertExpectations(t)
	})

}

func TestLogin(t *testing.T){
	t.Run("Login Success",func(t *testing.T){
		user := &models.Users{
			Email: "login@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		checkPassword := utils.NewHashMock()

		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		checkPassword.On("CheckPassword",user,user.Password).Return(true)
		
		userService := services.NewUserService(userRepo,checkPassword)

		token,returnUser,err := userService.Login(user)

		assert.NoError(t,err)
		assert.NotEmpty(t,token)
		assert.Equal(t, user.Email, returnUser.Email)
	})

	t.Run("Login password or email is required",func(t *testing.T) {
		user := &models.Users{
			Email: "",
			Password: "",
		}

		userRepo := repositories.NewUserRepositoryMock()
		checkPassword := utils.NewHashMock()
		
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		checkPassword.On("CheckPassword",user,user.Password).Return(false)

		userService := services.NewUserService(userRepo,checkPassword)

		_,_,err :=userService.Login(user)

		assert.EqualError(t,err,"Email or Password is required")

	})
}