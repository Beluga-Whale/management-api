package services

import (
	"errors"
	"fmt"

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

func (s *TaskService) GetAllTask(emailCookie string) ([]models.Tasks,error) {
	// NOTE - Decode Jwt in cookie เพื่อดึง Eamil
	email,err :=utils.ParseJWT(emailCookie)

	if err != nil{
		return nil,fmt.Errorf("Fail To Check Email : %w",err)
	}

	// NOTE - หา User จาก Email เพื่อเอา UserID 
	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return nil, errors.New("User not found")
	}

		return s.taskRepo.FindTaskAll(user.ID)
}

func (s *TaskService) FindTaskById(idSrt string, emailCookie string) (*models.Tasks, error) {

	// NOTE - Decode Jwt in cookie เพื่อดึง Eamil
	email,err :=utils.ParseJWT(emailCookie)

	if err != nil{
		return nil,fmt.Errorf("Fail To Check Email : %w",err)
	}

	// NOTE - หา User จาก Email เพื่อเอา UserID 
	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return nil, errors.New("User not found")
	}

	// NOTE -หา Task By ID
	task,err:= s.taskRepo.FindTaskById(idSrt)

	if err != nil {
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}

	// NOTE - มาเช็คว่าผู้ใช้เป็นเจ้าของ Task ไหม
	if task.UserID != user.ID {
		return nil, errors.New("you do not have permission to access this task")
	}
	
	return task,nil	

}


func (s*TaskService) UpdateTaskById(idStr string, emailCookie string, updatedTaskValue *models.Tasks) error {
	// NOTE - Check idStr
	if idStr == "" {
		return errors.New("Id is required")
	}

	// NOTE - Decode Jwt in cookie เพื่อดึง Eamil
	email,err :=utils.ParseJWT(emailCookie)

	if err != nil{
		return fmt.Errorf("Fail To Check Email : %w",err)
	}

	// NOTE - หา User จาก Email เพื่อเอา UserID 
	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return  errors.New("User not found")
	}


	// NOTE -หา Task By ID
	task,err:= s.taskRepo.FindTaskById(idStr)

	if err != nil {
		return  fmt.Errorf("failed to find task by ID: %w", err)
	}

	// NOTE - มาเช็คว่าผู้ใช้เป็นเจ้าของ Task ไหม
	if task.UserID != user.ID {
		return  errors.New("you do not have permission to access this task")
	}
	if	err :=s.taskRepo.UpdateTaskById(updatedTaskValue,task.ID); err != nil {
		return fmt.Errorf("Error :",err.Error())
	}
	
	return nil

}