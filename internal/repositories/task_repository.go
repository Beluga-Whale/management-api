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


func (repo *TaskRepository) FindTaskAll(userId uint, priority string) ([]models.Tasks,error) {
	var tasks []models.Tasks
	
	query := repo.db.Where("user_id = ?",userId)

	if priority != ""{
		query = query.Where("priority = ?",priority)
	}


	if err:= query.Order("created_at desc").Find(&tasks).Error; err !=nil {
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

func (repo *TaskRepository) DeleteTaskById(id uint) error {

	// ลบข้อมูลในฐานข้อมูลโดยใช้ id
	result := repo.db.Delete(&models.Tasks{}, id)
	if result.Error != nil {
		return result.Error // ส่งคืนข้อผิดพลาดหากการลบล้มเหลว
	}

	return nil // ส่งคืน nil หากลบสำเร็จ
}

func (repo *TaskRepository) FindTaskComplete(userId uint, priority string, complete bool) ([]models.Tasks,error) {
	var tasks []models.Tasks
	
	query := repo.db.Where("user_id = ?",userId).Where("completed = ?", complete)

	if priority != ""{
		query = query.Where("priority = ?",priority)
	}

	if err:= query.Order("created_at desc").Find(&tasks).Error; err !=nil {
		return  nil,err
	}
	return tasks,nil
}