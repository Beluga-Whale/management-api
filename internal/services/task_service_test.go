package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/Beluga-Whale/management-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T){
	t.Run("CreateTask Success",func(t *testing.T) {
		emailToken := "fakeToken"
		user:= &models.Users{
			Email: "Test@gmail.com",
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)

		userRepo.On("FindByEmail",user.Email).Return(user,nil)

		taskRepo.On("CreateTask",task).Return(nil)

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err :=taskService.CreateTask(task,emailToken)


		assert.NoError(t,err)
	})

	t.Run("Title and Description Required",func(t *testing.T) {
		emailToken := "fakeToken"
		task := &models.Tasks{
			Title: "",
			Description: "",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err := taskService.CreateTask(task,emailToken)

		assert.EqualError(t,err,"Title and Description is required")
	})

	t.Run("Jwt Error",func(t *testing.T) {
		emailToken := "fakeToken"

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("You not access to create task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err := taskService.CreateTask(task,emailToken)

		assert.EqualError(t,err,"Error JWT : You not access to create task")
	})

	t.Run("Error FindByEmail",func(t *testing.T) {
		emailToken := "fakeToken"
		user:= &models.Users{
			Email: "Test@gmail.com",
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("not have your email"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err := taskService.CreateTask(task,emailToken)

		assert.EqualError(t,err,"Not found user :not have your email")
	})

	t.Run("Create Error",func(t *testing.T) {
		emailToken := "fakeToken"
		user:= &models.Users{
			Email: "Test@gmail.com",
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)

		userRepo.On("FindByEmail",user.Email).Return(user,nil)

		taskRepo.On("CreateTask",task).Return(errors.New("You not create task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err :=taskService.CreateTask(task,emailToken)


		assert.EqualError(t,err,"You not create task")
	})
}

func TestGetAllTask(t *testing.T){
	t.Run("Get All Task Success",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskAll",user.ID,priority).Return([]models.Tasks{*task},nil)
		
		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		taskAll,err :=taskService.GetAllTask(emailToken,priority)

		assert.NoError(t,err)
		assert.Equal(t, []models.Tasks{*task}, taskAll)
		
		userRepo.AssertExpectations(t)
		jwtUtil.AssertExpectations(t)
	})

	t.Run("Fail To Check Email",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("check error"))


		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.GetAllTask(emailToken,priority)

		assert.EqualError(t,err,"Fail To Check Email : check error")
	})

	t.Run("User Not found",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("Can't to find you user"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err:=taskService.GetAllTask(emailToken,priority)

		assert.EqualError(t,err,"User not found")
	})
}

func TestFindTaskById(t *testing.T){
	t.Run("FindTaskById Success",func(t *testing.T) {
		idSrt := "1"
		emailToken := "fakeToken"

		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idSrt).Return(task,nil)

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		taskById,err:= taskService.FindTaskById(idSrt,emailToken)

		assert.NoError(t,err)
		assert.Equal(t,task,taskById)

		taskRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		jwtUtil.AssertExpectations(t)
				
	})

	t.Run("Fail To Check Email",func(t *testing.T) {
		idSrt := "1"
		emailToken := "fakeToken"

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("Can't to check your email"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.FindTaskById(idSrt,emailToken)
		assert.EqualError(t,err,"Fail To Check Email : Can't to check your email")
	})

	t.Run("User not found",func(t *testing.T) {
		idSrt := "1"
		emailToken := "fakeToken"

		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("User not found"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.FindTaskById(idSrt,emailToken)

		assert.EqualError(t,err,"User not found")
	})

	t.Run("Failed to find task by",func(t *testing.T) {
		idSrt := "1"
		emailToken := "fakeToken"

		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idSrt).Return(nil,errors.New("you can't to access this task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.FindTaskById(idSrt,emailToken)

		assert.EqualError(t,err,"failed to find task by ID: you can't to access this task")
	})

	t.Run("You not have permission to access this task",func(t *testing.T) {
		idSrt := "1"
		emailToken := "fakeToken"

		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 2,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idSrt).Return(task,nil)

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err:= taskService.FindTaskById(idSrt,emailToken)

		assert.EqualError(t,err,"you do not have permission to access this task")

		taskRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		jwtUtil.AssertExpectations(t)
	})
}

func TestUpdateTaskById(t *testing.T){
	t.Run("UpdateTaskById Success",func(t *testing.T) {
	idSrt := "1"
	emailToken := "fakeToken"

	user:= &models.Users{
		Email: "Test@gmail.com",
		Model: gorm.Model{ID: 1},
	}

	task := &models.Tasks{
		Title: "Title Test",
		Description: "Description Test",
		UserID: 1,
	}

	taskRepo := repositories.NewTaskRepositoryMock()
	userRepo := repositories.NewUserRepositoryMock()
	jwtUtil := utils.NewJwtMock()

	jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
	userRepo.On("FindByEmail",user.Email).Return(user,nil)
	taskRepo.On("FindTaskById",idSrt).Return(task,nil)
	taskRepo.On("UpdateTaskById",task,task.ID).Return(nil)

	taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

	err := taskService.UpdateTaskById(idSrt,emailToken,task)

	assert.NoError(t,err)

	userRepo.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
	jwtUtil.AssertExpectations(t)
	})

	t.Run("Id Is required",func(t *testing.T) {
		idStr := ""
		emailToken := "fakeToken"
			
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err :=taskService.UpdateTaskById(idStr,emailToken,task)

		assert.EqualError(t,err,"Id is required")

	})

	t.Run("Fail To Check Email",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("Not have your email"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.UpdateTaskById(idStr,emailToken,task)

		assert.EqualError(t,err,"Fail To Check Email : Not have your email")
	})

	t.Run("User not found",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("User not found"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.UpdateTaskById(idStr,emailToken,task)

		assert.EqualError(t,err,"User not found")
	})

	t.Run("Failed to find task by ID",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idStr).Return(nil,errors.New("Can't to find task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.UpdateTaskById(idStr,emailToken,task)

		assert.EqualError(t,err,"failed to find task by ID: Can't to find task")
	})

	t.Run("Not have permission to access this task",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 2,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idStr).Return(task,nil)

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.UpdateTaskById(idStr,emailToken,task)

		assert.EqualError(t,err,"you do not have permission to access this task")
	})

	t.Run("Not have permission to access this task",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idStr).Return(task,nil)
		taskRepo.On("UpdateTaskById",task,task.ID).Return(errors.New("Can't to update this task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.UpdateTaskById(idStr,emailToken,task)

		assert.EqualError(t,err,"Error : Can't to update this task")
	})
}

func TestDeleteTaskById(t *testing.T){
	t.Run("DeleteTaskById Success",func(t *testing.T) {
		idSrt := "1"
		emailToken := "fakeToken"

		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idSrt).Return(task,nil)
		taskRepo.On("DeleteTaskById",task.ID).Return(nil)
		
		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err := taskService.DeleteTaskById(idSrt,emailToken)

		assert.NoError(t,err)
		userRepo.AssertExpectations(t)
		taskRepo.AssertExpectations(t)
		jwtUtil.AssertExpectations(t)

	})

	t.Run("Id Is required",func(t *testing.T) {
		idStr := ""
		emailToken := "fakeToken"

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err :=taskService.DeleteTaskById(idStr,emailToken)

		assert.EqualError(t,err,"Id is required")

	})

	t.Run("Fail To Check Email",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("Not have your email"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.DeleteTaskById(idStr,emailToken)

		assert.EqualError(t,err,"Fail To Check Email : Not have your email")
	})

	t.Run("User not found",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("User not found"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.DeleteTaskById(idStr,emailToken)

		assert.EqualError(t,err,"User not found")
	})

	t.Run("Failed to find task by ID",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idStr).Return(nil,errors.New("Can't to find task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.DeleteTaskById(idStr,emailToken)

		assert.EqualError(t,err,"failed to find task by ID: Can't to find task")
	})

	
	t.Run("Not have permission to access this task",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 2,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idStr).Return(task,nil)

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.DeleteTaskById(idStr,emailToken)

		assert.EqualError(t,err,"you do not have permission to access this task")
	})

	t.Run("Fail To delete",func(t *testing.T) {
		idStr := "1"
		emailToken := "fakeToken"
			
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskById",idStr).Return(task,nil)
		taskRepo.On("DeleteTaskById",task.ID).Return(errors.New("You not delete this task"))

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		err:=taskService.DeleteTaskById(idStr,emailToken)

		assert.EqualError(t,err,"Error : You not delete this task")
	})
}

func TestGetCompleteTask(t *testing.T){
	t.Run("GetCompleteTask Success",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskComplete",user.ID,priority,true).Return([]models.Tasks{*task},nil)

		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		tasks,err :=taskService.GetCompleteTask(emailToken,priority)

		assert.NoError(t,err)
		assert.Equal(t,[]models.Tasks{*task},tasks)
	})

	t.Run("Fail to check email",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("Don't have your email"))


		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err :=taskService.GetCompleteTask(emailToken,priority)

		assert.EqualError(t,err,"Fail To Check Email : Don't have your email")
	})

	t.Run("User not found",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("User not found"))


		taskService := services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err :=taskService.GetCompleteTask(emailToken,priority)

		assert.EqualError(t,err,"User not found")
	})
}

func TestGetPendingTask(t *testing.T){
	t.Run("GetPendingTask Success",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskPending",user.ID,priority,true).Return([]models.Tasks{*task},nil)

		taskService:=services.NewTaskService(taskRepo,userRepo,jwtUtil)

		tasks,err := taskService.GetPendingTask(emailToken,priority)

		assert.NoError(t,err)

		assert.Equal(t,[]models.Tasks{*task},tasks)

	})

	t.Run("Fail To Check Email",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("Not have email"))


		taskService:=services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.GetPendingTask(emailToken,priority)


		assert.EqualError(t,err,"Fail To Check Email : Not have email")

	})

	t.Run("GetPendingTask Success",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("User not found"))


		taskService:=services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.GetPendingTask(emailToken,priority)


		assert.EqualError(t,err,"User not found")

	})
}

func TestGetOverdueTask(t *testing.T){
	t.Run("GetOverdueTask Success",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Test",
			UserID: 1,
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(user,nil)
		taskRepo.On("FindTaskOverdue",user.ID,priority,true).Return([]models.Tasks{*task},nil)

		taskService:=services.NewTaskService(taskRepo,userRepo,jwtUtil)

		tasks,err := taskService.GetOverdueTask(emailToken,priority)

		assert.NoError(t,err)

		assert.Equal(t,[]models.Tasks{*task},tasks)

	})

	t.Run("Fail To Check Email",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("",errors.New("Not have email"))


		taskService:=services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.GetOverdueTask(emailToken,priority)


		assert.EqualError(t,err,"Fail To Check Email : Not have email")

	})

	t.Run("GetOverdueTask Fail",func(t *testing.T) {
		emailToken := "fakeToken"
		priority := ""
		user:= &models.Users{
			Email: "Test@gmail.com",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		userRepo := repositories.NewUserRepositoryMock()
		jwtUtil := utils.NewJwtMock()

		jwtUtil.On("ParseJWT",emailToken).Return("Test@gmail.com",nil)
		userRepo.On("FindByEmail",user.Email).Return(nil,errors.New("User not found"))


		taskService:=services.NewTaskService(taskRepo,userRepo,jwtUtil)

		_,err := taskService.GetOverdueTask(emailToken,priority)


		assert.EqualError(t,err,"User not found")

	})
}