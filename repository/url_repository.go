package repository

import (
	"url-shortener/models"

	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(url models.URL) error
	FindByUrl(url string) (models.URL, error)
	FindById(id string) (models.URL, error)
	Update(url models.URL) error
	Delete(url models.URL) error
}

type urlRepository struct {
	db *gorm.DB
}

func (u *urlRepository) Create(url models.URL) error {
	//TODO implement me
	panic("implement me")
}

func (u *urlRepository) FindByUrl(url string) (models.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (u *urlRepository) FindById(id string) (models.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (u *urlRepository) Update(url models.URL) error {
	//TODO implement me
	panic("implement me")
}

func (u *urlRepository) Delete(url models.URL) error {
	//TODO implement me
	panic("implement me")
}

func NewUrlRepository(db *gorm.DB) UrlRepository {
	return &urlRepository{
		db: db,
	}
}
