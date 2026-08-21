package service

import (
	"context"
	"errors"
	"net/http"
	"time"
	"url-shortener/appErrors"
	"url-shortener/dto"
	"url-shortener/models"
	"url-shortener/repository"
	"url-shortener/utils"

	"github.com/redis/go-redis/v9"
)

type UrlService interface {
	CreateUrl(userId uint, req dto.CreateUrlRequest) (*models.URL, error)
	GetUrl(userId uint, urlId uint) (*models.URL, error)
	GetUserUrls(userId uint) (*[]models.URL, error)
	UpdateUrl(userId uint, urlId uint, req dto.UpdateUrlRequest) (*models.URL, error)
	DeleteUrl(userId uint, urlId uint) error
	ActivateUrl(userId uint, urlId uint) error
	DeactivateUrl(userId uint, urlId uint) error
	Redirect(shortCode string, ctx context.Context) (string, error)
}

type urlService struct {
	urlRepo repository.UrlRepository
	redis   *redis.Client
}

func NewUrlService(urlRepo repository.UrlRepository, redis *redis.Client) UrlService {
	return &urlService{
		urlRepo: urlRepo,
		redis:   redis,
	}
}

func (u *urlService) CreateUrl(userId uint, req dto.CreateUrlRequest) (*models.URL, error) {
	if req.CustomAlias != nil {
		var exists bool
		exists, _ = u.urlRepo.ExistsByShortCode(*req.CustomAlias)

		if exists {
			return nil, &appErrors.AppError{
				Code:    "URl_ALREADY_EXISTS",
				Message: "Url already exists",
				Status:  http.StatusConflict,
			}
		}
	}

	shortCode, err := utils.GenerateShortCode()
	if err != nil {
		return nil, err
	}

	url := &models.URL{
		OriginalURL: req.OriginalUrl,
		ShortCode:   shortCode,
		UserID:      userId,
		ExpiresAt:   req.ExpiresAt,
		IsActive:    true,
		IsDeleted:   false,
	}

	if err := u.urlRepo.Create(url); err != nil {
		return nil, err
	}

	return url, nil

}

func (u *urlService) GetUrl(userId uint, urlId uint) (*models.URL, error) {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (u *urlService) GetUserUrls(userId uint) (*[]models.URL, error) {
	urls, err := u.urlRepo.FindAllByUserId(userId)

	if err != nil {
		return nil, &appErrors.AppError{
			Code:    "DB_ERROR",
			Message: "Database error",
			Status:  http.StatusInternalServerError,
		}
	}

	return urls, nil
}

func (u *urlService) UpdateUrl(userId uint, urlId uint, req dto.UpdateUrlRequest) (*models.URL, error) {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)

	if err != nil {
		return nil, appErrors.ErrURLNotFound
	}

	url.OriginalURL = req.OriginalUrl

	if req.ExpiresAt != nil {
		url.ExpiresAt = req.ExpiresAt
	}

	if req.IsActive != nil {
		url.IsActive = *req.IsActive
	}

	ok := u.urlRepo.Update(url)

	if ok != nil {
		return nil, &appErrors.AppError{
			Code:    "DB_ERROR",
			Message: "Database error",
			Status:  http.StatusInternalServerError,
		}
	}

	return url, nil
}

func (u *urlService) DeleteUrl(userId uint, urlId uint) error {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)
	if err != nil {
		return appErrors.ErrURLNotFound
	}

	if err := u.urlRepo.SoftDelete(url); err != nil {
		return appErrors.ErrServerError
	}
	return nil
}

func (u *urlService) ActivateUrl(userId uint, urlId uint) error {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)

	if err != nil {
		return appErrors.ErrURLNotFound
	}

	url.IsActive = true

	if err := u.urlRepo.Update(url); err != nil {
		return appErrors.ErrServerError
	}

	return nil
}

func (u *urlService) DeactivateUrl(userId uint, urlId uint) error {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)

	if err != nil {
		return appErrors.ErrURLNotFound
	}

	url.IsActive = false

	if err := u.urlRepo.Update(url); err != nil {
		return appErrors.ErrServerError
	}

	return nil
}

func (u *urlService) Redirect(shortCode string, ctx context.Context) (string, error) {

	key := "url:" + shortCode

	originalUrl, err := u.redis.Get(ctx, key).Result()

	if err == nil {
		// Cache HIT
		return originalUrl, nil
	}

	if !errors.Is(err, redis.Nil) {
		return "", err
	}

	// Cache MISS
	url, err := u.urlRepo.FindByShortCode(shortCode)

	if err != nil {
		return "", appErrors.ErrURLNotFound
	}

	if url.IsDeleted || !url.IsActive {
		return "", &appErrors.AppError{
			Code:    "URL_ALREADY_DELETED",
			Message: "Url already deleted",
			Status:  http.StatusNoContent,
		}
	}

	if url.ExpiresAt != nil && !url.ExpiresAt.After(time.Now()) {
		url.IsActive = false
		if err := u.urlRepo.UpdateIsActive(url.ID, false); err != nil {
			return "", appErrors.ErrServerError
		}
		return "", &appErrors.AppError{
			Code:    "URL_EXPIRES_ALREADY",
			Message: "Url already expired",
			Status:  http.StatusNoContent,
		}
	}

	err = u.redis.Set(
		ctx,
		key,
		url.OriginalURL,
		10*time.Minute,
	).Err()

	if err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}
