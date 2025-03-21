package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Beluga-Whale/management-api/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository{
	return &TaskRepository{db:db}
}

func (repo *TaskRepository) CreateTask(task *models.Tasks)error {
	if strings.Trim(task.Title,"") == ""{
		return errors.New("Task Title can't be empty")
	}

	return repo.db.Create(task).Error
}


func (repo *TaskRepository) FindTaskAll(userId uint) ([]models.Tasks,error) {

	var tasks []models.Tasks
	
	if err:= repo.db.Where("user_id",userId).Find(&tasks).Error; err !=nil {
		return  nil,err
	}
	return tasks,nil
}

func (repo *TaskRepository) FindTaskById(idStr string) (*models.Tasks, error) {
	var task models.Tasks

	id,err :=  strconv.Atoi(idStr)

	if err !=nil{
		return nil,errors.New("Invalid Task ID fomat")
	}

	result := repo.db.First(&task, id)

	if result.Error != nil {
		return nil, result.Error
	}
	
	return &task,nil
}

func (repo *TaskRepository) UpdateTaskById(updatedTaskValue *models.Tasks, taskID uint) error {
	var task models.Tasks

	result := repo.db.First(&task, taskID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("Can't find task by ID: %w", result.Error)
		}
		return fmt.Errorf("Error fetching task: %w", result.Error)
	}

	// NOTE - Update task 
	if err:= repo.db.Model(&task).Updates(updatedTaskValue).Error; err != nil {
		return fmt.Errorf("Failed to update task: %w",err)
	} 

	return nil
}