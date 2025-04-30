package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Beluga-Whale/management-api/config"
	"github.com/Beluga-Whale/management-api/internal/handlers"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/Beluga-Whale/management-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setUpAppUser() *fiber.App {
	os.Setenv("HOST", "localhost")      
	os.Setenv("PORT", "5433")       
	os.Setenv("DATABASE_NAME", "taskManage_test") 
	os.Setenv("USER_NAME", "postgres")   
	os.Setenv("PASSWORD", "password") 
	config.ConnectTestDB()
	hashUtil := utils.NewHash()
	jwtUtil := utils.NewJwt()

	userRepo := repositories.NewUserRepository(config.TestDB)
	userService := services.NewUserService(userRepo, hashUtil, jwtUtil)
	userHandler := handlers.NewUserHandler(userService)



	app := fiber.New()

	app.Post("/user/register", userHandler.RegisterUser)
	app.Post("/user/login", userHandler.Login)
	app.Put("/user/:id", userHandler.EditUser)

	return app
}


func clearDataBaseUser(){
	if err := config.TestDB.Exec("DELETE FROM tasks").Error; err != nil {
        log.Fatalf("Failed to clear tasks table: %v", err)
    }

	if err := config.TestDB.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clear test database: %v", err)
	}
	
}

func TestUserRegisterIntegration(t *testing.T){
	t.Run("Test User Registration Failed", func(t *testing.T) {
		app := setUpAppUser()
		req := httptest.NewRequest("POST", "/user/register",nil)
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid request")
		clearDataBaseUser()
	})
	t.Run("Test User Registration Success", func(t *testing.T) {
		app := setUpAppUser()

		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		reqBody := []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusCreated,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User registered success")
		clearDataBaseUser()
	})

	t.Run("Test User Registration Success Already Exist", func(t *testing.T) {
		app := setUpAppUser()
	
		//NOTE - สร้างอีเมลที่ไม่ซ้ำกัน
		email := fmt.Sprintf("test_integration_%d@gmail.com", time.Now().Unix())
		reqBody := []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))
	
		req := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
	
		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
	
		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User registered success")
	
		//NOTE - ทดสอบการลงทะเบียนด้วยอีเมลที่ซ้ำกัน
		reqBody = []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))
		req = httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
	
		res, err = app.Test(req)
	
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	
		body, _ = io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Email has already been used")
		clearDataBaseUser()
	})


}

func TestLoginIntegration(t *testing.T){
	t.Run("Test Login Error Invalid",func(t *testing.T) {
		app := setUpAppUser()

		req := httptest.NewRequest("POST", "/user/login",nil)
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)	
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"In valid request")
		clearDataBaseUser()
	})

	t.Run("Test Login Success",func(t *testing.T) {

		app := setUpAppUser()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		reqBody := []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusCreated,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User registered success")


		req = httptest.NewRequest("POST","/user/login", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie","jwt=fake-jwt-token")

		res,err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		body,_ = io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Login success")
		clearDataBaseUser()
	})

	t.Run("Test Login Invalid Password or Email",func(t *testing.T) {

		app := setUpAppUser()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		reqBody := []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))

		req := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusCreated,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User registered success")


		req = httptest.NewRequest("POST","/user/login", bytes.NewReader([]byte(`{"email":"wrongemail","password":"wrongpassword"}`)))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie","jwt=fake-jwt-token")

		res,err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		body,_ = io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid Email or Password")
		clearDataBaseUser()
	})
}

func TestEditUserIntegration(t *testing.T){
	t.Run("Test Edit User Error",func(t *testing.T) {
		app := setUpAppUser()

		req := httptest.NewRequest("PUT", "/user/1",nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie","jwt=fake-jwt-token")

		res,err := app.Test(req)

		assert.NoError(t, err)	
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid request")
		clearDataBaseUser()
	})

	t.Run("Test Edit User Success",func(t *testing.T) {
		app := setUpAppUser()

		// NOTE - Register User
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		reqBody := []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))
		req := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User registered success")

		// NOTE - Login User
		req = httptest.NewRequest("POST", "/user/login", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fake-jwt") 
		res, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		loginBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(loginBody), "Login success")

		// Get the JWT token from the login response
		var loginResponse map[string]interface{}
		if err := json.Unmarshal(loginBody, &loginResponse); err != nil {
			t.Fatal("Error parsing login response:", err)
		}

		jwtToken, ok := loginResponse["token"].(string)
		if !ok {
			t.Fatal("Expected token in response, but got:", loginResponse)
		}

		// ดึง user id จากการตอบกลับ
		user, ok := loginResponse["user"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected user in response, but got:", loginResponse)
		}
		userID, ok := user["id"].(float64)
		if !ok {
			t.Fatal("Expected user ID in response, but got:", loginResponse)
		}

		// NOTE - Update User
		updatedUser := `{"name":"Updated User","bio":"Updated bio"}`
		req = httptest.NewRequest("PUT", fmt.Sprintf("/user/%d", int(userID)), bytes.NewReader([]byte(updatedUser)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=" + jwtToken)
		res, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		//NOTE - ตรวจสอบว่าได้ข้อมูลที่แก้ไข
		body, _ = io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Updated User")
		assert.Contains(t, string(body), "Updated bio")
		clearDataBaseUser()
	})
}

