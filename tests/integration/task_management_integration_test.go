package integration_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Beluga-Whale/management-api/config"
	"github.com/Beluga-Whale/management-api/internal/handlers"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/Beluga-Whale/management-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setUpAppTask() *fiber.App {
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "5433")
	os.Setenv("DATABASE_NAME", "taskManage_test")
	os.Setenv("USER_NAME", "postgres")
	os.Setenv("PASSWORD", "password")
	config.ConnectTestDB()
	jwtUtil := utils.NewJwt()
	hashUtil := utils.NewHash()

	userRepo := repositories.NewUserRepository(config.TestDB)
	taskRepo := repositories.NewTaskRepository(config.TestDB)

	taskService := services.NewTaskService(taskRepo,userRepo, jwtUtil)
	userService := services.NewUserService(userRepo, hashUtil, jwtUtil)

	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	app := fiber.New()

	app.Post("/user/register", userHandler.RegisterUser)

	app.Get("/task", taskHandler.GetAllTask)
	app.Post("/task", taskHandler.CreateTask)
	app.Get("/task/complete", taskHandler.GetCompleteTask)
	app.Get("/task/pending", taskHandler.GetPendingTask)
	app.Get("/task/overdue", taskHandler.GetOverdueTask)
	return app
}

func clearDataBaseTask(){
	if err := config.TestDB.Exec("DELETE FROM tasks").Error; err != nil {
		log.Fatalf("Failed to clear test tasks database: %v", err)
	}
	if err := config.TestDB.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clear test users database: %v", err)
	}
}


// NOTE - Create Fucntion 
func createJWT(email string) (string, error) {
	jwtUtil := utils.NewJwt()
	return jwtUtil.GenerateJWT(email)
}

func registerUser(t *testing.T, email string){
	app := setUpAppTask()
	// NOTE - Create User
	reqBody := []byte(fmt.Sprintf(`{"email":"%s","password":"password1234","name":"Test User"}`, email))

	req := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	res,err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t,fiber.StatusCreated,res.StatusCode)

	body,_ := io.ReadAll(res.Body)
	assert.Contains(t,string(body),"User registered success")
}

func registerAndCreateTask(t *testing.T,email string,token string){
	app := setUpAppTask()

	// NOTE - Create User
	registerUser(t,email)
	// NOTE - Create TASK
	token, err := createJWT(email)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	reqBody := []byte(fmt.Sprintf(`{"title":"Test Task","description":"This is a test task"}`))
	req := httptest.NewRequest("POST", "/task", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "jwt=" + token) // ใส่ JWT ใน Cookie

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(body), "create task success")
}

// NOTE - Integration Test
func TestCreateTaskIntegration(t *testing.T){
	
	t.Run("Create Task Success", func(t *testing.T) {
        app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())

		// NOTE - Create User
		registerUser(t,email)
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

        reqBody := []byte(fmt.Sprintf(`{"title":"Test Task","description":"This is a test task"}`))
        req := httptest.NewRequest("POST", "/task", bytes.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Cookie", "jwt=" + token) // ใส่ JWT ใน Cookie

        res, err := app.Test(req)
        assert.NoError(t, err)
        assert.Equal(t, fiber.StatusOK, res.StatusCode)

        body, _ := io.ReadAll(res.Body)
        assert.Contains(t, string(body), "create task success")

        clearDataBaseTask() 
    })

	t.Run("Create Task Error", func(t *testing.T) {
        app := setUpAppTask()

		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		registerUser(t,email)

		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

  
        req := httptest.NewRequest("POST", "/task", nil)
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Cookie", "jwt=" + token) // ใส่ JWT ใน Cookie

        res, err := app.Test(req)
        assert.NoError(t, err)
        assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

        body, _ := io.ReadAll(res.Body)
        assert.Contains(t, string(body), "Invalid request")

        clearDataBaseTask() 
    })
}

func TestGetAllTaskIntegration(t *testing.T){
	t.Run("GetAllTask Success",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())

		// NOTE - Create User
		registerUser(t,email)

        token, err := createJWT(email)
        if err != nil {
			t.Fatalf("Failed to create JWT: %v", err)
        }

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task", nil)
        req.Header.Set("Cookie", "jwt=" + token) //NOTE - ใส่ JWT ใน Cookie

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		clearDataBaseTask()
	})

	t.Run("GetAllTask Fail to check Email",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())

		// NOTE - Create User
		registerUser(t,email)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task", nil)
        req.Header.Set("Cookie", "jwt=fsdfsd") //NOTE - ใส่ JWT ใน Cookie

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Fail To Check Email")
		clearDataBaseTask() 
	})

}

func TestGetCompleteTask(t *testing.T){
	t.Run("GetCompleteTask Success",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/complete", nil)
		req.Header.Set("Cookie", "jwt=" + token)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		clearDataBaseTask()
	})

	t.Run("GetCompleteTask User not authenticated",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/complete", nil)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusUnauthorized,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User not authenticated")
		clearDataBaseTask()
	})
	
	t.Run("GetCompleteTask User Fail to get all tasks",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/complete", nil)
		req.Header.Set("Cookie", "jwt=invalid-jwt" + token)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Fail To Check Email")
		clearDataBaseTask()
	})
}
func TestGetPendingTask(t *testing.T){
	t.Run("GetPendingTask Success",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/pending", nil)
		req.Header.Set("Cookie", "jwt=" + token)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		clearDataBaseTask()
	})

	t.Run("GetPendingTask User not authenticated",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/pending", nil)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusUnauthorized,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User not authenticated")
		clearDataBaseTask()
	})
	
	t.Run("GetPendingTask User Fail to get all tasks",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/pending", nil)
		req.Header.Set("Cookie", "jwt=invalid-jwt" + token)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Fail To Check Email")
		clearDataBaseTask()
	})
}

func TestGetOverdueTask(t *testing.T){
	t.Run("GetOverdueTask Success",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/overdue", nil)
		req.Header.Set("Cookie", "jwt=" + token)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		clearDataBaseTask()
	})

	t.Run("GetOverdueTask User not authenticated",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/overdue", nil)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusUnauthorized,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User not authenticated")
		clearDataBaseTask()
	})
	
	t.Run("GetOverdueTask User Fail to get all tasks",func(t *testing.T) {
		app := setUpAppTask()
		email := fmt.Sprintf("test_integration_%s@gmail.com", uuid.NewString())
		
		// NOTE - Create TASK
        token, err := createJWT(email)
        if err != nil {
            t.Fatalf("Failed to create JWT: %v", err)
        }

		registerAndCreateTask(t,email,token)

		// NOTE - Get All TASK
		req := httptest.NewRequest("GET", "/task/overdue", nil)
		req.Header.Set("Cookie", "jwt=invalid-jwt" + token)
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Fail To Check Email")
		clearDataBaseTask()
	})
}