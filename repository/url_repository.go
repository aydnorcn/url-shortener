package repository

import (
	"url-shortener/models"

	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(url *models.URL) error
	FindById(id uint) (*models.URL, error)
	FindByIdAndUserId(id uint, userId uint) (*models.URL, error)
	FindByShortCode(code string) (*models.URL, error)
	ExistsByShortCode(code string) (bool, error)
	FindAllByUserId(userId uint) (*[]models.URL, error)
	CountByUserId(userId uint) (int64, error)
	Update(url *models.URL) error
	SoftDelete(url *models.URL) error
	SetActive(url *models.URL) error
}

type urlRepository struct {
	db *gorm.DB
}

func (u *urlRepository) ExistsByShortCode(code string) (bool, error) {
	var exists bool

	err := u.db.Raw(`
        SELECT EXISTS (
            SELECT 1 FROM urls WHERE short_code = ?
        )
    `, code).Scan(&exists).Error

	return exists, err
}

func NewUrlRepository(db *gorm.DB) UrlRepository {
	return &urlRepository{
		db: db,
	}
}

func (u *urlRepository) Create(url *models.URL) error {
	return u.db.Create(url).Error
}

func (u *urlRepository) FindById(id uint) (*models.URL, error) {
	var url models.URL
	err := u.db.Where("id = ?", id).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (u *urlRepository) FindByIdAndUserId(id uint, userId uint) (*models.URL, error) {
	var url models.URL
	err := u.db.Where("id = ? AND user_id = ?", id, userId).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (u *urlRepository) FindByShortCode(code string) (*models.URL, error) {
	var url models.URL
	err := u.db.Where("short_code = ?", code).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (u *urlRepository) FindAllByUserId(userId uint) (*[]models.URL, error) {
	var urls []models.URL

	err := u.db.Where("user_id = ?", userId).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return &urls, nil
}

func (u *urlRepository) CountByUserId(userId uint) (int64, error) {
	var count int64
	err := u.db.Model(&models.URL{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *urlRepository) Update(url *models.URL) error {
	var urlModel models.URL
	err := u.db.Model(&urlModel).Where("id = ?", url.ID).Updates(url).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *urlRepository) SoftDelete(url *models.URL) error {
	url.IsDeleted = true
	url.IsActive = false
	return u.db.Save(url).Error
}

func (u *urlRepository) SetActive(url *models.URL) error {
	url.IsActive = true
	return u.db.Save(url).Error
}
