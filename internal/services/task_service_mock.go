package services

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type TaskServiceMock struct {
	mock.Mock
}

func NewTaskServiceMock() *TaskServiceMock{
	return &TaskServiceMock{}
}

func (m *TaskServiceMock) CreateTask(task *models.Tasks, emailCookie string) error {
	args := m.Called(task,emailCookie)
	return args.Error(0)
}

func (m *TaskServiceMock) GetAllTask(emailCookie string ,priority string) ([]models.Tasks,error) {
	args := m.Called(emailCookie, priority)
	return  args.Get(0).([]models.Tasks), args.Error(1)
}

func  (m *TaskServiceMock) FindTaskById(idSrt string, emailCookie string) (*models.Tasks, error) {
	args := m.Called(idSrt,emailCookie)

	return args.Get(0).(*models.Tasks),args.Error(1)
}

func (m *TaskServiceMock) UpdateTaskById(idStr string, emailCookie string, updatedTaskValue *models.Tasks) error  {
	args := m.Called(idStr,emailCookie,updatedTaskValue)

	return args.Error(0)
}

func (m *TaskServiceMock) DeleteTaskById(idStr string,emailCookie string,) error   {
	args := m.Called(idStr,emailCookie)

	return args.Error(0)
}

func (m *TaskServiceMock) GetCompleteTask(emailCookie string ,priority string) ([]models.Tasks,error)   {
	args := m.Called(emailCookie,priority)

	return args.Get(0).([]models.Tasks),args.Error(1)
}

func (m *TaskServiceMock) GetPendingTask(emailCookie string ,priority string) ([]models.Tasks,error) {
	args := m.Called(emailCookie,priority)

	return args.Get(0).([]models.Tasks),args.Error(1)
}

func (m *TaskServiceMock) GetOverdueTask(emailCookie string ,priority string) ([]models.Tasks,error) {
	args := m.Called(emailCookie,priority)

	return args.Get(0).([]models.Tasks),args.Error(1)
}


