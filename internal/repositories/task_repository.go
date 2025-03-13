package repositories

import (
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
	return repo.db.Create(task).Error
}

// func (repo *TaskRepository) GetTaskByUserId(userID uint) error {
// 	 err:= repo.db.Where("userID")
// }