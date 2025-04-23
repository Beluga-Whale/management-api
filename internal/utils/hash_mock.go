package utils

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type HashMock struct {
	mock.Mock
}

func NewHashMock() *HashMock {
	return &HashMock{}
}

func (m *HashMock) CheckPassword(user *models.Users, password string) bool {
	args :=m.Called(user,password)
	return args.Bool(0)
}
