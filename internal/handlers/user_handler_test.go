package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/management-api/internal/handlers"
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	t.Run("Test Register Success",func(t *testing.T) {
		// NOTE - Arrange
		user := &models.Users{
			Email: "Test@gmail.com",
			Password: "password123",
			Name: "Test User",
		}
		
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		userService.On("RegisterUser",user).Return(nil)

		app := fiber.New()
		app.Post("/user/register",userHandler.RegisterUser)

		reqBody := []byte(`{
			"email":"Test@gmail.com",
			"password":"password123",
			"name":"Test User"
		}`)

		req := httptest.NewRequest("POST","/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User registered success")

		userService.AssertExpectations(t)
	})

	t.Run("RegisterUser ฺBody empty Invalid request", func(t *testing.T) {
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Post("/register", userHandler.RegisterUser)

		req := httptest.NewRequest("POST", "/register", nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request")
	})

	t.Run("Test Register Error",func(t *testing.T) {
		// NOTE - Arrange
		user := &models.Users{
			Email: "Test@gmail.com",
			Password: "password123",
			Name: "Test User",
		}
		
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		userService.On("RegisterUser",user).Return(errors.New("User already exists"))

		app := fiber.New()
		app.Post("/user/register",userHandler.RegisterUser)

		reqBody := []byte(`{
			"email":"Test@gmail.com",
			"password":"password123",
			"name":"Test User"
		}`)

		req := httptest.NewRequest("POST","/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User already exists")

		userService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T){
	t.Run("Test Login Success",func(t *testing.T) {
		userLogin := &models.Users{
			Email:"test@gmail.com",
			Password:"password123",
		}

		jwtToken := "fake_jwtToken"

		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		expectUser := &models.Users{
			Email: "test@gmail.com",
			Name: "Test User",
		}
		userService.On("Login",userLogin).Return(jwtToken,expectUser,nil)	

		app:= fiber.New()
		app.Post("/user/login",userHandler.Login)

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password123"
		}`)
		req := httptest.NewRequest("POST","/user/login",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)

		body,_ := io.ReadAll(res.Body)

		assert.Contains(t,string(body),"Login success")
	})
	t.Run("Test Login BadRequest",func(t *testing.T) {
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		app:= fiber.New()
		app.Post("/user/login",userHandler.Login)

		
		req := httptest.NewRequest("POST","/user/login",nil)
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body,_ := io.ReadAll(res.Body)

		assert.Contains(t,string(body),"In valid request")
	})

	t.Run("Test Login Error",func(t *testing.T) {
		userLogin := &models.Users{
			Email:"test@gmail.com",
			Password:"password123",
		}
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		userService.On("Login",userLogin).Return("",nil,errors.New("User not found"))	

		app:= fiber.New()
		app.Post("/user/login",userHandler.Login)

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password123"
		}`)
		req := httptest.NewRequest("POST","/user/login",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body,_ := io.ReadAll(res.Body)

		assert.Contains(t,string(body),"User not found")
	})
}

func TestLogout(t *testing.T){
	t.Run("Test Logout Success",func(t *testing.T) {
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Post("/user/logout",userHandler.Logout)


		req := httptest.NewRequest("POST","/user/logout",nil)

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)

		body,_ := io.ReadAll(res.Body)

		assert.Contains(t,string(body),"Logout Success")
	})
}

func TestGetUser(t *testing.T){
	t.Run("GetUser Success", func(t *testing.T) {
		userEmail := "test@example.com"
		expectedUser := &models.Users{
			Email: "test@example.com",
			Name:  "Test User",
		}

		userService := new(services.UserServiceMock)
		userHandler := handlers.NewUserHandler(userService)

		userService.On("GetUserByEmail", userEmail).Return(expectedUser, nil)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userEmail", userEmail)
			return c.Next()
		}

		app := fiber.New()
		app.Get("/user", testMiddleware,userHandler.GetUser)

		req := httptest.NewRequest("GET", "/user", nil)
		
		res, err := app.Test(req)
	
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Test User")
		assert.Contains(t, string(body), "test@example.com")

		userService.AssertExpectations(t)
	})

	t.Run("GetUser ฺBadRequest", func(t *testing.T) {
		userEmail := "test@example.com"

		userService := new(services.UserServiceMock)
		userHandler := handlers.NewUserHandler(userService)

		userService.On("GetUserByEmail", userEmail).Return(nil, errors.New("User not found"))

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userEmail", userEmail)
			return c.Next()
		}

		app := fiber.New()
		app.Get("/user", testMiddleware,userHandler.GetUser)

		req := httptest.NewRequest("GET", "/user", nil)
		
		res, err := app.Test(req)
	
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User not found")

		userService.AssertExpectations(t)
	})
}

func TestEditUser(t *testing.T){
	t.Run("EditUser Success",func(t *testing.T) {
		userMock := &models.Users{
			Name: "Edit User",
			Email: "test@gmail.com",
		}

		emailCookie := "fake_jwtToken"

		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		userService.On("UpdateUserById","1",emailCookie,userMock).Return(nil)

		app := fiber.New()
		app.Put("/user/:id", userHandler.EditUser)

		reqBody := []byte(`{
			"name":"Edit User",
			"email":"test@gmail.com"
		}`)

		req := httptest.NewRequest("PUT", "/user/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+emailCookie)

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

	})

	t.Run("EditUser ID is required",func(t *testing.T) {

		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Put("/user", userHandler.EditUser)


		req := httptest.NewRequest("PUT", "/user",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Task ID is required")
	})
	
	t.Run("EditUser not authenticated",func(t *testing.T) {
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)


		app := fiber.New()
		app.Put("/user/:id", userHandler.EditUser)

		reqBody := []byte(`{
			"name":"Edit User",
			"email":"test@gmail.com"
		}`)

		req := httptest.NewRequest("PUT", "/user/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)	

		assert.Contains(t, string(body), "User not authenticated")
	})

	t.Run("EditUser Invalid request",func(t *testing.T) {
		userMock := &models.Users{
			Name: "Edit User",
			Email: "test@gmail.com",
		}

		emailCookie := "fake_jwtToken"

		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		userService.On("UpdateUserById","1",emailCookie,userMock).Return(nil)

		app := fiber.New()
		app.Put("/user/:id", userHandler.EditUser)

		req := httptest.NewRequest("PUT", "/user/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+emailCookie)

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request")
	})

	t.Run("EditUser Success",func(t *testing.T) {
		userMock := &models.Users{
			Name: "Edit User",
			Email: "test@gmail.com",
		}

		emailCookie := "fake_jwtToken"

		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)

		userService.On("UpdateUserById","1",emailCookie,userMock).Return(errors.New("User not found"))

		app := fiber.New()
		app.Put("/user/:id", userHandler.EditUser)

		reqBody := []byte(`{
			"name":"Edit User",
			"email":"test@gmail.com"
		}`)

		req := httptest.NewRequest("PUT", "/user/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+emailCookie)

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)	
		assert.Contains(t, string(body), "User not found")

	})
}