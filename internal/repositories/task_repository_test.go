package repositories_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T){
	t.Run("CreateTask success",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Test Title",
			Description: "Test Description",
		}

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("CreateTask",task).Return(nil)

		err :=taskRepo.CreateTask(task)

		assert.NoError(t,err)
		taskRepo.AssertExpectations(t)
	})

	t.Run("CreateTask Error",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Test Title",
			Description: "Test Description",
		}

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("CreateTask",task).Return(errors.New("Error create task"))

		err :=taskRepo.CreateTask(task)

		assert.EqualError(t,err,"Error create task")
		taskRepo.AssertExpectations(t)
	})
	
	
}

func TestFindTaskAll(t *testing.T){
	t.Run("FindTaskAll success",func(t *testing.T) {
		priority := ""
		user := &models.Users{
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "TEST",
			Description: "test",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("FindTaskAll",user.ID,priority).Return([]models.Tasks{*task},nil)

	
		tasks,err:=taskRepo.FindTaskAll(user.ID,priority)

		assert.NoError(t,err)
		assert.Equal(t, []models.Tasks{*task},tasks)
	})
	t.Run("Error FindTaskAll",func(t *testing.T) {
		priority := ""
		user := &models.Users{
			Model: gorm.Model{ID: 1},
		}


		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("FindTaskAll",user.ID,priority).Return(nil,errors.New("Can not find this task"))

	
		_,err:=taskRepo.FindTaskAll(user.ID,priority)

		
		assert.EqualError(t, err,"Can not find this task")
	})
}

func TestFindTaskById(t *testing.T){
	t.Run("FindTaskById success",func(t *testing.T) {
		idStr := ""

		task := &models.Tasks{
			Title: "TEST",
			Description: "test",
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("FindTaskById",idStr).Return(task,nil)

	
		tasks,err:=taskRepo.FindTaskById(idStr)

		assert.NoError(t,err)
		assert.Equal(t, task,tasks)
	})
	t.Run("Error FindTaskById",func(t *testing.T) {
		idStr := ""

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("FindTaskById",idStr).Return(nil,errors.New("Can't find task by your id"))

		_,err:=taskRepo.FindTaskById(idStr)

		assert.EqualError(t,err,"Can't find task by your id")
	})
}

func TestUpdateTaskById(t *testing.T){
	t.Run("UpdateTaskById Success",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Test Title",
			Description: "Description Test",
			Model: gorm.Model{ID: 1},
		}
		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("UpdateTaskById",task,task.ID).Return(nil)

		err :=taskRepo.UpdateTaskById(task,task.ID)

		assert.NoError(t,err)
		taskRepo.AssertExpectations(t)

	})

	t.Run("Error UpdateTaskById ",func(t *testing.T) {
		task := &models.Tasks{
			Title: "Test Title",
			Description: "Description Test",
			Model: gorm.Model{ID: 1},
		}
		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("UpdateTaskById",task,task.ID).Return(errors.New("Can't update this task"))

		err :=taskRepo.UpdateTaskById(task,task.ID)

		assert.EqualError(t,err,"Can't update this task")

	})
}

func TestDeleteTaskById(t *testing.T){
	t.Run("Delete success",func(t *testing.T) {

		taskRepo := repositories.NewTaskRepositoryMock()
		taskRepo.On("DeleteTaskById",uint(1)).Return(nil)

		err :=taskRepo.DeleteTaskById(uint(1))
		 assert.NoError(t,err)
		 taskRepo.AssertExpectations(t)
	})
	t.Run("Error Delete ",func(t *testing.T) {

		taskRepo := repositories.NewTaskRepositoryMock()
		taskRepo.On("DeleteTaskById",uint(1)).Return(errors.New("Can't delete this task"))

		err :=taskRepo.DeleteTaskById(uint(1))
		assert.EqualError(t,err,"Can't delete this task")

	})
}

func TestFindTaskComplete(t *testing.T){
	t.Run("FindTaskComplete success",func(t *testing.T) {
		priority := ""
		user := &models.Users{
			Model: gorm.Model{ID: 1},
		}

		task := &models.Tasks{
			Title: "Title Test",
			Description: "Description Tets",
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		taskRepo.On("FindTaskComplete",user.ID,priority,true).Return([]models.Tasks{*task},nil)

		tasks,err :=taskRepo.FindTaskComplete(user.ID,priority,true)
		 assert.NoError(t,err)
		 assert.Equal(t,[]models.Tasks{*task},tasks)
		 taskRepo.AssertExpectations(t)
	})
	t.Run("Error FindTaskComplete",func(t *testing.T) {
		priority := ""
		user := &models.Users{
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()
		taskRepo.On("FindTaskComplete",user.ID,priority,true).Return(nil,errors.New("Can't to find task complete"))

		_,err :=taskRepo.FindTaskComplete(user.ID,priority,true)
		 assert.EqualError(t,err,"Can't to find task complete")
		 taskRepo.AssertExpectations(t)
	})
}

func TestFindTaskPending(t *testing.T){
	t.Run("FindTaskPending Success",func(t *testing.T) {
		priority := ""
		user := &models.Users{
			Model: gorm.Model{ID: 1},
		}
		task := &models.Tasks{
			Title: "test",
		}

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("FindTaskPending",user.ID,priority,false).Return([]models.Tasks{*task},nil)

		tasks,err :=taskRepo.FindTaskPending(user.ID,priority,false)

		assert.NoError(t,err)
		assert.Equal(t,[]models.Tasks{*task},tasks)
		taskRepo.AssertExpectations(t)
	})
	t.Run("Error FindTaskPending ",func(t *testing.T) {
		priority := ""
		user := &models.Users{
			Model: gorm.Model{ID: 1},
		}

		taskRepo := repositories.NewTaskRepositoryMock()

		taskRepo.On("FindTaskPending",user.ID,priority,false).Return(nil,errors.New("Can't to find this task"))

		_,err :=taskRepo.FindTaskPending(user.ID,priority,false)


		assert.EqualError(t,err,"FindTaskPending")
		taskRepo.AssertExpectations(t)
	})
}

func TestFindTaskOverdue(t *testing.T){
	t.Run("FindTaskOverdue Success",func(t *testing.T) {
		priority :=""
	task := &models.Tasks{
		Title: "Test title",
		Description: "Test description",
		Model: gorm.Model{ID:1},
	}
	user := &models.Users{
		Model: gorm.Model{ID: 1},
	}

	taskRepo := repositories.NewTaskRepositoryMock()
	
	taskRepo.On("FindTaskOverdue",user.ID,priority,true).Return([]models.Tasks{*task},nil)

	tasks,err:=taskRepo.FindTaskOverdue(user.ID,priority,true)

	assert.NoError(t,err)
	assert.Equal(t,[]models.Tasks{*task},tasks)
	})
	t.Run("Error FindTaskOverdue",func(t *testing.T) {
		priority :=""

	user := &models.Users{
		Model: gorm.Model{ID: 1},
	}

	taskRepo := repositories.NewTaskRepositoryMock()
	
	taskRepo.On("FindTaskOverdue",user.ID,priority,true).Return(nil,errors.New("Can't find task overdue"))

	_,err:=taskRepo.FindTaskOverdue(user.ID,priority,true)

	assert.EqualError(t,err,"Can't find task overdue");
	})
}