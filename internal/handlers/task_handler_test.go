package handlers_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/management-api/internal/handlers"
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllTask(t *testing.T) {
	t.Run("GetAllTask Success", func(t *testing.T) {
		task := []models.Tasks{
			{Title: "Test Task",
			Description: "This is a test task",},
			{Title: "Test Task2",
			Description: "This is a test task2",},
		}

		emailJwt := "fakeJWT@gmail.com"
		priority := "high"

		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetAllTask",emailJwt,priority).Return(task,nil)

		app := fiber.New()
		app.Get("/tasks",taskHandler.GetAllTask)

		// NOTE - httptest.NewRequest ส่งจะ 3 ตัว 1. method 2. url 3. body
		req := httptest.NewRequest("GET",fmt.Sprintf("/tasks?priority=%s",priority),nil)
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err := app.Test(req)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusOK,res.StatusCode)

		body,_ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Test Task")
		assert.Contains(t,string(body),"Test Task2")
		taskService.AssertExpectations(t)
	})
	t.Run("GetAllTask Not authenticated", func(t *testing.T) {
		emailJwt := ""
		priority := "high"

		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetAllTask",emailJwt,priority).Return(nil,errors.New("User not authenticated"))

		app := fiber.New()
		app.Get("/tasks",taskHandler.GetAllTask)

		// NOTE - httptest.NewRequest ส่งจะ 3 ตัว 1. method 2. url 3. body
		req := httptest.NewRequest("GET",fmt.Sprintf("/tasks?priority=%s",priority),nil)

		res,_ := app.Test(req)

		assert.Equal(t,fiber.StatusUnauthorized,res.StatusCode)

	})

	t.Run("GetAllTask BadRequest", func(t *testing.T) {

		emailJwt := "fakeJWT@gmail.com"
		priority := "high"

		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetAllTask",emailJwt,priority).Return(nil,errors.New("Can't get all tasks"))

		app := fiber.New()
		app.Get("/tasks",taskHandler.GetAllTask)

		// NOTE - httptest.NewRequest ส่งจะ 3 ตัว 1. method 2. url 3. body
		req := httptest.NewRequest("GET",fmt.Sprintf("/tasks?priority=%s",priority),nil)
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err := app.Test(req)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Can't get all tasks")
		taskService.AssertExpectations(t)
	})
}

func TestCreateTask(t *testing.T){
	t.Run("CreateTask Success", func(t *testing.T) {
		// NOTE - Arrange
		task := &models.Tasks{
			Title:       "Test Task",
			Description: "This is a test task",
		}
		emailJwt := "fakeJWT@gmail.com"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("CreateTask", task, emailJwt).Return(nil)

		app := fiber.New()
		app.Post("/task", taskHandler.CreateTask)

		reqBody := []byte(`{
			"title": "Test Task",
			"description": "This is a test task"
		}`)
		req := httptest.NewRequest("POST", "/task", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		// NOTE - Assert 
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		taskService.AssertExpectations(t)
	})

	t.Run("CreateTask BodyParser BadRequest", func(t *testing.T) {
		// NOTE -  Arrange
		task := &models.Tasks{
			Title: 	 "Test Task",
			Description: "This is a test task",
		}
		
		emailJwt := "fake@gmail.com"
		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("CreateTask", task, emailJwt).Return(nil)

		app := fiber.New()
		app.Post("/task", taskHandler.CreateTask)

		req := httptest.NewRequest("POST", "/task", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert 
		assert.NoError(t, err)
		
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Invalid request")

	})

	t.Run("CreateTask User not Authenticated", func(t *testing.T) {
		// NOTE - Arrange
		task := &models.Tasks{
			Title:       "Test Task",
			Description: "This is a test task",
		}
		emailJwt := ""

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("CreateTask", task, emailJwt).Return(nil)

		app := fiber.New()
		app.Post("/task", taskHandler.CreateTask)

		reqBody := []byte(`{
			"title": "Test Task",
			"description": "This is a test task"
		}`)
		req := httptest.NewRequest("POST", "/task", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// NOTE - Act 
		res, err := app.Test(req)

		// NOTE - Assert 
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	})

	t.Run("CreateTask Fail", func(t *testing.T) {
		// NOTE - Arrange
		task := &models.Tasks{
			Title:       "Test Task",
			Description: "This is a test task",
		}
		emailJwt := "fakeJWT@gmail.com"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("CreateTask", task, emailJwt).Return(errors.New("Can't create task"))

		app := fiber.New()
		app.Post("/task", taskHandler.CreateTask)

		reqBody := []byte(`{
			"title": "Test Task",
			"description": "This is a test task"
		}`)
		req := httptest.NewRequest("POST", "/task", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		// NOTE - Assert 
		assert.NoError(t, err)

		body,_ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Can't create task")
	})
}

func TestFindTaskById(t *testing.T){
	t.Run("FindTaskById Success", func(t *testing.T) {
		// NOTE - Arrange
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
		}
		idStr:= "1"

		emailCookie := "fakeJWT@gmail.com"

		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("FindTaskById", idStr,emailCookie).Return(task,nil)

		app := fiber.New()
		app.Get("/task/:id", taskHandler.FindTaskById)

		req := httptest.NewRequest("GET", "/task/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusOK, res.StatusCode)
		assert.Contains(t, string(body), "Title Test")
	})

	t.Run("FindTaskById Unauthenticated", func(t *testing.T) {
		// NOTE - Arrange
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		app := fiber.New()
		app.Get("/task/:id", taskHandler.FindTaskById)

		req := httptest.NewRequest("GET", "/task/1", nil)
		req.Header.Set("Cookie", "jwt=")
		// NOTE - Act 
		res, err := app.Test(req)

		// NOTE - Assert
		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User not authenticated")
	})

	t.Run("FindTaskById ID is required", func(t *testing.T) {
		// NOTE - Arrange
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)
		
		emailCookie := "fakeJWT@gmail.com"
		taskService.On("FindTaskById", "1",emailCookie).Return(nil,errors.New("User not authenticated"))

		app := fiber.New()
		app.Get("/task", taskHandler.FindTaskById)

		req := httptest.NewRequest("GET", "/task", nil)
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")
		// NOTE - Act 
		res, err := app.Test(req)

		// NOTE - Assert
		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "sk ID is required")
	})

	t.Run("FindTaskById Success", func(t *testing.T) {
		// NOTE - Arrange
		idStr:= "1"

		emailCookie := "fakeJWT@gmail.com"

		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("FindTaskById", idStr,emailCookie).Return(nil,errors.New("Can't find task"))

		app := fiber.New()
		app.Get("/task/:id", taskHandler.FindTaskById)

		req := httptest.NewRequest("GET", "/task/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Can't find task")
	})
}

func TestUpdateTask(t *testing.T){
	t.Run("UpdateTask Success",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Updated Task",
			Description: "This is an updated task",
		}
		emailCookie := "fakeJWT@gmail.com"
		idStr := "1"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("UpdateTaskById",idStr,emailCookie,task).Return(nil)

		app:= fiber.New()
		app.Put("/task/:id",taskHandler.UpdateTask)

		reqBody := []byte(`{
		"Title": "Updated Task",
			"Description": "This is an updated task"
		}`)
		req := httptest.NewRequest("PUT","/task/1",bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type","application/json")
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		assert.Contains(t,string(body),"Update Task Success")
	})

	t.Run("UpdateTask ID is required",func(t *testing.T) {
		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		app:= fiber.New()
		app.Put("/task/",taskHandler.UpdateTask)

		req := httptest.NewRequest("PUT","/task/",nil)

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		assert.Contains(t,string(body),"Task ID is required")
	})

	t.Run("UpdateTask User not authenticated",func(t *testing.T) {
		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		app:= fiber.New()
		app.Put("/task/:id",taskHandler.UpdateTask)

		req := httptest.NewRequest("PUT","/task/1",nil)

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusUnauthorized,res.StatusCode)
		assert.Contains(t,string(body),"User not authenticated")
	})

	t.Run("UpdateTask Invalid",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Updated Task",
			Description: "This is an updated task",
		}
		emailCookie := "fakeJWT@gmail.com"
		idStr := "1"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("UpdateTaskById",idStr,emailCookie,task).Return(nil)

		app:= fiber.New()
		app.Put("/task/:id",taskHandler.UpdateTask)

		req := httptest.NewRequest("PUT","/task/1",nil)
		req.Header.Set("Content-Type","application/json")
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		assert.Contains(t,string(body),"Invalid request")
	})

	t.Run("UpdateTask Error",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Updated Task",
			Description: "This is an updated task",
		}
		emailCookie := "fakeJWT@gmail.com"
		idStr := "1"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("UpdateTaskById",idStr,emailCookie,task).Return(errors.New("Can't update task"))

		app:= fiber.New()
		app.Put("/task/:id",taskHandler.UpdateTask)

		reqBody := []byte(`{
		"Title": "Updated Task",
			"Description": "This is an updated task"
		}`)
		req := httptest.NewRequest("PUT","/task/1",bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type","application/json")
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		assert.Contains(t,string(body),"Can't update task")
	})
}

func TestDelete(t *testing.T){
	t.Run("DeleteTask Success",func(t *testing.T) {

		emailCookie := "fakeJWT@gmail.com"
		idStr := "1"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("DeleteTaskById",idStr,emailCookie).Return(nil)

		app:= fiber.New()
		app.Delete("/task/:id",taskHandler.DeleteTask)

		req := httptest.NewRequest("DELETE","/task/1",nil)
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		assert.Contains(t,string(body),"Delete Task Success")
	})

	t.Run("DeleteTask ID is required",func(t *testing.T) {

		emailCookie := "fakeJWT@gmail.com"
		idStr := ""

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("DeleteTaskById",idStr,emailCookie).Return(nil)

		app:= fiber.New()
		app.Delete("/task",taskHandler.DeleteTask)

		req := httptest.NewRequest("DELETE","/task",nil)
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		assert.Contains(t,string(body),"Task ID is required")
	})
	t.Run("DeleteTask Not authenticated",func(t *testing.T) {

		emailCookie := ""
		idStr := "1"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("DeleteTaskById",idStr,emailCookie).Return(nil)

		app:= fiber.New()
		app.Delete("/task/:id",taskHandler.DeleteTask)

		req := httptest.NewRequest("DELETE","/task/1",nil)
		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusUnauthorized,res.StatusCode)
		assert.Contains(t,string(body),"User not authenticated")
	})

	t.Run("DeleteTask BadRequest",func(t *testing.T) {

		emailCookie := "fakeJWT@gmail.com"
		idStr := "1"

		taskService := new(services.TaskServiceMock)
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("DeleteTaskById",idStr,emailCookie).Return(errors.New("Can't delete task"))

		app:= fiber.New()
		app.Delete("/task/:id",taskHandler.DeleteTask)

		req := httptest.NewRequest("DELETE","/task/1",nil)
		req.Header.Set("Cookie","jwt=fakeJWT@gmail.com")

		res,err:=app.Test(req)

		body,_:= io.ReadAll(res.Body)

		assert.NoError(t,err)

		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		assert.Contains(t,string(body),"Can't delete task")
	})
}

func TestGetCompleteTask(t *testing.T) {
	t.Run("GetCompleteTask Success",func(t *testing.T) {
		task := []models.Tasks{
			{Title: "Complete Task1",
			Description: "This is a complete task",},
			{Title: "Complete Task2",
			Description: "This is a complete task2",},
		}

		// emailCookie := "fakeJWT@gmail.com"
		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetCompleteTask",mock.Anything,mock.Anything).Return(task,nil)

		app := fiber.New()
		app.Get("/task/complete",taskHandler.GetCompleteTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/complete?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusOK, res.StatusCode)
		assert.Contains(t, string(body), "Complete Task1")
		assert.Contains(t, string(body), "Complete Task2")

		taskService.AssertExpectations(t)
	})

	t.Run("GetCompleteTask not authenticated",func(t *testing.T) {
		task := []models.Tasks{
			{Title: "Complete Task1",
			Description: "This is a complete task",},
			{Title: "Complete Task2",
			Description: "This is a complete task2",},
		}

		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetCompleteTask",mock.Anything,mock.Anything).Return(task,nil)

		app := fiber.New()
		app.Get("/task/complete",taskHandler.GetCompleteTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/complete?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, string(body), "User not authenticated")
	})

	t.Run("GetCompleteTask not authenticated",func(t *testing.T) {

		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetCompleteTask",mock.Anything,mock.Anything).Return(nil,errors.New("Error to get complete task"))

		app := fiber.New()
		app.Get("/task/complete",taskHandler.GetCompleteTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/complete?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Error to get complete task")
	})
}

func TestGetPendingTask(t *testing.T) {
	t.Run("GetPendingTask Success",func(t *testing.T) {
		task := []models.Tasks{
			{Title: "Pending Task1",
			Description: "This is a pending task",},
			{Title: "Pending Task2",
			Description: "This is a pending task2",},
		}

		// emailCookie := "fakeJWT@gmail.com"
		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetPendingTask",mock.Anything,mock.Anything).Return(task,nil)

		app := fiber.New()
		app.Get("/task/pending",taskHandler.GetPendingTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/pending?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusOK, res.StatusCode)
		assert.Contains(t, string(body), "Pending Task1")
		assert.Contains(t, string(body), "Pending Task2")

		taskService.AssertExpectations(t)
	})

	t.Run("GetPendingTask not authenticated",func(t *testing.T) {
		task := []models.Tasks{
			{Title: "Pending Task1",
			Description: "This is a pending task",},
			{Title: "Pending Task2",
			Description: "This is a pending task2",},
		}

		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetPendingTask",mock.Anything,mock.Anything).Return(task,nil)

		app := fiber.New()
		app.Get("/task/pending",taskHandler.GetPendingTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/pending?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, string(body), "User not authenticated")
	})

	t.Run("GetPendingTask Error to get pending task",func(t *testing.T) {

		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetPendingTask",mock.Anything,mock.Anything).Return(nil,errors.New("Error to get pending task"))

		app := fiber.New()
		app.Get("/task/pending",taskHandler.GetPendingTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/pending?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Error to get pending task")
	})
}

func TestGetOverdueTask(t *testing.T) {
	t.Run("GetOverdueTask Success",func(t *testing.T) {
		task := []models.Tasks{
			{Title: "OverdueTask Task1",
			Description: "This is a overdueTask task",},
			{Title: "OverdueTask Task2",
			Description: "This is a overdueTask task2",},
		}

		// emailCookie := "fakeJWT@gmail.com"
		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetOverdueTask",mock.Anything,mock.Anything).Return(task,nil)

		app := fiber.New()
		app.Get("/task/overdueTask",taskHandler.GetOverdueTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/overdueTask?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusOK, res.StatusCode)
		assert.Contains(t, string(body), "OverdueTask Task1")
		assert.Contains(t, string(body), "OverdueTask Task2")

		taskService.AssertExpectations(t)
	})

	t.Run("GetOverdueTask not authenticated",func(t *testing.T) {
		task := []models.Tasks{
			{Title: "OverdueTask Task1",
			Description: "This is a overdueTask task",},
			{Title: "OverdueTask Task2",
			Description: "This is a overdueTask task2",},
		}

		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetOverdueTask",mock.Anything,mock.Anything).Return(task,nil)

		app := fiber.New()
		app.Get("/task/overdueTask",taskHandler.GetOverdueTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/overdueTask?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, string(body), "User not authenticated")
	})

	t.Run("GetOverdueTask Error to get overdueTask task",func(t *testing.T) {

		priority := "high"
		taskService := services.NewTaskServiceMock()
		taskHandler := handlers.NewTaskHandler(taskService)

		taskService.On("GetOverdueTask",mock.Anything,mock.Anything).Return(nil,errors.New("Error to get overdueTask task"))

		app := fiber.New()
		app.Get("/task/overdueTask",taskHandler.GetOverdueTask)

		req := httptest.NewRequest("GET",fmt.Sprintf("/task/overdueTask?priority=%s",priority),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=fakeJWT@gmail.com")

		// NOTE - Act 
		res, err := app.Test(req)

		body,_ := io.ReadAll(res.Body)

		// NOTE - Assert
		assert.NoError(t, err)

		assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Error to get overdueTask task")
	})
}