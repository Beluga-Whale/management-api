package repositories

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository{
	return &UserRepository{db:db}
}

func (repo *UserRepository) CreateUser(user *models.Users)error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return repo.db.Create(user).Error
}

func (repo *UserRepository) FindByEmail(email string) (*models.Users, error ){
	var user models.Users
	
	result := repo.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound{
			return nil,nil
		}
		return nil, result.Error
	}

	return &user, nil

}
