package services

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (m *UserServiceMock) RegisterUser(user *models.Users) error {
	args :=m.Called(user)
	return args.Error(0)
}

func (m *UserServiceMock) Login(user *models.Users) (string,*models.Users,error) {
	args :=m.Called(user)
	if task,ok := args.Get(1).(*models.Users) ; ok {
		return args.String(0),task,nil
	}
	return "",nil,args.Error(2)
}

func (m *UserServiceMock) GetUserByEmail(email string) (*models.Users, error) {
	args :=m.Called(email)
	if user,ok := args.Get(0).(*models.Users);ok{
		return user,nil
	}
	return nil,args.Error(1)
}

func (m *UserServiceMock) UpdateUserById(idStr string, emailCookie string, updatedUserValue *models.Users) (error) {
	args :=m.Called(idStr,emailCookie,updatedUserValue)
	return args.Error(0)
}

