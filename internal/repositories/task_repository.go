package repositories

import (
	"errors"
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


func (repo *TaskRepository) FindTaskAll() ([]models.Tasks,error) {
	var tasks []models.Tasks
	if err:= repo.db.Find(&tasks).Error; err !=nil {
		return  nil,err
	}
	return tasks,nil

}
