package services

import (
	"errors"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/repositories"
	"github.com/Beluga-Whale/management-api/internal/utils"
)

type TaskService struct {
	taskRepo *repositories.TaskRepository
	userRepo *repositories.UserRepository
}

func NewTaskService(taskRepo *repositories.TaskRepository, userRepo *repositories.UserRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, userRepo:userRepo}
}

func (s *TaskService)  CreateTask(task *models.Tasks, emailCookie string) error {
	if task.Title =="" || task.Description == ""{
		return errors.New("Title and Description is required")
	}

	claims, err := utils.ParseJWT(emailCookie)

	if err != nil {
		return errors.New("Error JWT : "+err.Error())
	}

	user, err := s.userRepo.FindByEmail( claims)

	if err != nil {
		return errors.New("Not found user :"+err.Error())
	}

	task.UserID = user.ID

	if err :=s.taskRepo.CreateTask(task); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) GetAllTask() ([]models.Tasks,error) {
	return s.taskRepo.FindTaskAll()
}