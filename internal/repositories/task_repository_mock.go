package repositories

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type TaskRepositoryMock struct {
	mock.Mock
}

func NewTaskRepositoryMock() *TaskRepositoryMock{
	return &TaskRepositoryMock{}
}

func (m *TaskRepositoryMock) CreateTask(task *models.Tasks)error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *TaskRepositoryMock) FindTaskAll(userId uint, priority string) ([]models.Tasks, error) {
	args := m.Called(userId, priority)
	return args.Get(0).([]models.Tasks), args.Error(1)
}

func (m *TaskRepositoryMock) FindTaskById(idStr string) (*models.Tasks, error) {
	args := m.Called(idStr)
	if task, ok := args.Get(0).(*models.Tasks); ok {
		return task, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TaskRepositoryMock) UpdateTaskById(updatedTaskValue *models.Tasks, taskID uint) error {
	args := m.Called(updatedTaskValue,taskID)
	return args.Error(0)
}

func (m *TaskRepositoryMock) DeleteTaskById(id uint) error{
	args := m.Called(id)
	return args.Error(0)
}

func (m *TaskRepositoryMock)FindTaskComplete(userId uint, priority string, complete bool) ([]models.Tasks,error){
	args := m.Called(userId,priority,complete)

	if tasks, ok := args.Get(0).([]models.Tasks); ok {
		return tasks, nil
	}

	return nil, args.Error(1)
}

func (m *TaskRepositoryMock)FindTaskPending(userId uint, priority string, complete bool) ([]models.Tasks,error){
	args := m.Called(userId,priority,complete)
	return args.Get(0).([]models.Tasks),args.Error(1)
}

func (m *TaskRepositoryMock)FindTaskOverdue(userId uint, priority string, complete bool) ([]models.Tasks,error){
	args := m.Called(userId,priority,complete)
	return args.Get(0).([]models.Tasks),args.Error(1)
}



