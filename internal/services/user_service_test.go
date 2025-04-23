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
	"gorm.io/gorm"
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
		
		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

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
		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)
	
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

		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

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

		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

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

		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

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
		hashUtil := utils.NewHashMock()

		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		hashUtil.On("CheckPassword",user,user.Password).Return(true)
		
		jwtUtil := utils.NewJwtMock()
		jwtUtil.On("GenerateJWT",user.Email).Return("token",nil)
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

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
		hashUtil := utils.NewHashMock()
		
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		hashUtil.On("CheckPassword",user,user.Password).Return(false)

		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

		_,_,err :=userService.Login(user)

		assert.EqualError(t,err,"Email or Password is required")

	})

	t.Run("Login Email not found",func(t *testing.T) {
		user := &models.Users{
			Email: "test@",
			Password: "1234551",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashUtil := utils.NewHashMock()

		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("Error fail"))
		
		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

		_,_,err := userService.Login(user)

		assert.EqualError(t,err,"Email not found")
	})

	t.Run("Login Invalid Email or Password",func(t *testing.T) {
		user := &models.Users{
			Email: "test@",
			Password: "1234551",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashUtil := utils.NewHashMock()

		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		hashUtil.On("CheckPassword",user,user.Password).Return(false)

		jwtUtil := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

		_,_,err := userService.Login(user)

		assert.EqualError(t,err,"Invalid Email or Password")
		userRepo.AssertExpectations(t)
	})

	t.Run("Login Fail to generate token",func(t *testing.T) {
		user := &models.Users{
			Email: "login@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashUtil := utils.NewHashMock()

		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		hashUtil.On("CheckPassword",user,user.Password).Return(true)
		
		jwtUtil := utils.NewJwtMock()
		jwtUtil.On("GenerateJWT",user.Email).Return("",errors.New("Failed to generate JWT"))
		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

		_,_,err := userService.Login(user)

		assert.EqualError(t,err,"Failed to generate token: Failed to generate JWT")
	})
}

func TestGetUserByEmail(t *testing.T){
	t.Run("GetUser Success",func(t *testing.T){
		user := &models.Users{
			Email: "login@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("FindByEmail",user.Email).Return(user,nil)

		userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

		user,err := userService.GetUserByEmail(user.Email)

		assert.NoError(t,err)
		assert.NotEmpty(t,user)
		userRepo.AssertExpectations(t)
	})


}

func TestUpdateUserById(t *testing.T){
	// t.Run("Update Success",func(t *testing.T){

	// 	idUser := "1"
	// 	emailToken := "fakeToken"
	// 	// NOTE User
	// 	user := &models.Users{
	// 		Email: "already@mail.com",
	// 		Password: "1222123",
	// 		Name: "tester",
	// 		Model: gorm.Model{
	// 			ID: 1,
	// 		},
	// 	}
	// 	// NOTE - Update Bio
	// 	userUpdate := &models.Users{
	// 		Bio: "Test Update Bio",
	// 	}
	
	// 	userRepo := repositories.NewUserRepositoryMock()
	// 	hashUtil := utils.NewHashMock()
	// 	jwtUtil := utils.NewJwtMock()

	// 	jwtUtil.On("ParseJWT",emailToken).Return("test@gmail.com",nil)
	// 	userRepo.On("FindByEmail","test@gmail.com").Return(user,nil)

	// 	userRepo.On("FindUserById",idUser).Return(user,nil)

	// 	userRepo.On("UpdateUserById",idUser).Return(nil)
		
	// 	userService := services.NewUserService(userRepo,hashUtil,jwtUtil)

	// 	err := userService.UpdateUserById(idUser,emailToken,userUpdate)

	// 	assert.NoError(t,err)
	// })

	t.Run("Update Success", func(t *testing.T) {
		idUser := "1"
		emailToken := "fakeToken"
		user := &models.Users{
			Bio:    "Test Update Bio",
			Email:  "test@gmail.com",
			Model:  gorm.Model{ID: 1}, // กำหนด ID เป็น 1
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashUtil := utils.NewHashMock()
		jwtUtil := utils.NewJwtMock()

		// Mock ParseJWT ให้คืน email ที่ถูกต้อง
		jwtUtil.On("ParseJWT", emailToken).Return("test@gmail.com", nil)

		// Mock FindByEmail ให้คืน user จริงๆ
		userRepo.On("FindByEmail", user.Email).Return(user, nil)

		// Mock FindUserById ให้คืน user ที่มี ID ตรงกับที่ต้องการ
		userRepo.On("FindUserById", idUser).Return(user, nil)

		// Mock UpdateUserById ให้คืน nil (ไม่มีข้อผิดพลาด)
		userRepo.On("UpdateUserById", user, user.ID).Return(nil)

		userService := services.NewUserService(userRepo, hashUtil, jwtUtil)

		// เรียก UpdateUserById
		err := userService.UpdateUserById(idUser, emailToken, user)

		// ตรวจสอบว่าไม่มีข้อผิดพลาด
		assert.NoError(t, err)

		// ตรวจสอบว่า mock ถูกเรียกตามที่คาดหวัง
		jwtUtil.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})
}
