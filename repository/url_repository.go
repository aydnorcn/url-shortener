package repository

import (
	"url-shortener/models"

	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(url models.URL) error
	FindByUrl(url string) (*models.URL, error)
	FindById(id string) (*models.URL, error)
	Update(url *models.URL) error
	Delete(url *models.URL) error
}

type urlRepository struct {
	db *gorm.DB
}

func (u *urlRepository) Create(url models.URL) error {
	return u.db.Create(url).Error
}

func (u *urlRepository) FindByUrl(url string) (*models.URL, error) {
	var urlModel models.URL
	err := u.db.Where("url = ?", url).First(&urlModel).Error
	if err != nil {
		return nil, err
	}
	return &urlModel, nil
}

func (u *urlRepository) FindById(id string) (*models.URL, error) {
	var url models.URL
	err := u.db.Where("id = ?", id).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (u *urlRepository) Update(url *models.URL) error {
	var urlModel models.URL
	err := u.db.Model(&urlModel).Where("id = ?", url.ID).Updates(url).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *urlRepository) Delete(url *models.URL) error {
	return u.db.Delete(&url).Error
}

func NewUrlRepository(db *gorm.DB) UrlRepository {
	return &urlRepository{
		db: db,
	}
}
