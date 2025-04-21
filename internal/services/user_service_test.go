package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/services"
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
		

		userService := services.NewUserService(userRepo)

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
		userService := services.NewUserService(userRepo)
	
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

		userService := services.NewUserService(userRepo)

		err := userService.RegisterUser(user)

		assert.EqualError(t,err,"Fail To Check Email : Fail To Check Email")
	})
}