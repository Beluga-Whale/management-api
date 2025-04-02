package repositories

import (
	"errors"
	"fmt"
	"strconv"

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

func (repo *UserRepository) FindUserById(idStr string) (*models.Users, error ){
	var user models.Users
	id,err :=  strconv.Atoi(idStr)

	if err !=nil{
		return nil,errors.New("Invalid Task ID fomat")
	}

	result := repo.db.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}
	
	return &user,nil
}

func (repo *UserRepository) UpdateUserById(updatedUserValue *models.Users, userID uint) error {
	var user models.Users

	result := repo.db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("Can't find user by ID: %w", result.Error)
		}
		return fmt.Errorf("Error fetching user: %w", result.Error)
	}

	// NOTE - Update task 
	if err:= repo.db.Model(&user).Updates(updatedUserValue).Error; err != nil {
		return fmt.Errorf("Failed to update user: %w",err)
	} 

	return nil
}
