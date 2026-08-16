package service

import (
	"url-shortener/dto"
	"url-shortener/models"
	"url-shortener/repository"
)

type UrlService interface {
	CreateUrl(userId uint, req dto.CreateUrlRequest) (*models.URL, error)
	GetUrl(userId, urlId uint) (*models.URL, error)
	UpdateUrl(userId uint, urlId uint, req dto.UpdateUrlRequest) (*models.URL, error)
	DeleteUrl(userId uint, urlId uint) error
	Redirect(shortCode string, userId uint) error
}

type urlService struct {
	urlRepo repository.UrlRepository
}

func NewUrlService(urlRepo repository.UrlRepository) UrlService {
	return &urlService{
		urlRepo: urlRepo,
	}
}

func (u *urlService) CreateUrl(userId uint, req dto.CreateUrlRequest) (*models.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (u *urlService) GetUrl(userId, urlId uint) (*models.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (u *urlService) UpdateUrl(userId uint, urlId uint, req dto.UpdateUrlRequest) (*models.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (u *urlService) DeleteUrl(userId uint, urlId uint) error {
	//TODO implement me
	panic("implement me")
}

func (u *urlService) Redirect(shortCode string, userId uint) error {
	//TODO implement me
	panic("implement me")
}
