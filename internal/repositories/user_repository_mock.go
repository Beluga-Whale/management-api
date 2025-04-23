package repositories

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock{
	return &UserRepositoryMock{}
}

func (m *UserRepositoryMock) CreateUser(user *models.Users)error{
	args :=m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindByEmail(email string) (*models.Users, error ){
	args :=m.Called(email)

	// NOTE- ต้องเช็คว่า args.Get(0) เป็น nil ไหม ไม่งั้นมันจะ panic
	if user, ok := args.Get(0).(*models.Users); ok {
        return user, args.Error(1)
    }

	return nil, args.Error(1)
}

func (m *UserRepositoryMock)FindUserById(idStr string) (*models.Users, error ){
	args :=m.Called(idStr)

	if user,ok := args.Get(0).(*models.Users); ok{
		return user, args.Error(1)
	}

	return  nil,args.Error(1)
}

func (m *UserRepositoryMock) UpdateUserById(updatedUserValue *models.Users, userID uint) error{
	args :=m.Called(updatedUserValue,userID)
	return args.Error(0)
}

