package repository

import (
	"url-shortener/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindById(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(user *models.User) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindById(id uint) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}
