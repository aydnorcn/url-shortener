package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"url-shortener/appErrors"
	"url-shortener/cache"
	"url-shortener/dto"
	"url-shortener/metrics"
	"url-shortener/models"
	"url-shortener/repository"
	"url-shortener/utils"
)

type CachedURL struct {
	ID          uint       `json:"id"`
	OriginalURL string     `json:"original_url"`
	IsActive    bool       `json:"is_active"`
	IsDeleted   bool       `json:"is_deleted"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type UrlService interface {
	CreateUrl(userId uint, req dto.CreateUrlRequest) (*models.URL, error)
	GetUrl(userId uint, urlId uint) (*models.URL, error)
	GetUserUrls(userId uint, page, pageSize int) (*[]models.URL, int64, error)
	UpdateUrl(userId uint, urlId uint, req dto.UpdateUrlRequest) (*models.URL, error)
	DeleteUrl(userId uint, urlId uint) error
	ActivateUrl(userId uint, urlId uint) error
	DeactivateUrl(userId uint, urlId uint) error
	Redirect(shortCode string, ctx context.Context) (*models.URL, error)
}

type urlService struct {
	urlRepo repository.UrlRepository
	cache   cache.Cache
}

func NewUrlService(urlRepo repository.UrlRepository, cache cache.Cache) UrlService {
	return &urlService{
		urlRepo: urlRepo,
		cache:   cache,
	}
}

func (u *urlService) CreateUrl(userId uint, req dto.CreateUrlRequest) (*models.URL, error) {
	var shortCode string
	if req.CustomAlias != nil && *req.CustomAlias != "" {
		exists, _ := u.urlRepo.ExistsByShortCode(*req.CustomAlias)
		if exists {
			return nil, &appErrors.AppError{
				Code:    "URL_ALREADY_EXISTS",
				Message: "Url alias already exists",
				Status:  http.StatusConflict,
			}
		}
		shortCode = *req.CustomAlias
	} else {
		generated, err := utils.GenerateShortCode()
		if err != nil {
			return nil, err
		}
		shortCode = generated
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
		return nil, appErrors.ErrURLNotFound
	}
	return url, nil
}

func (u *urlService) GetUserUrls(userId uint, page, pageSize int) (*[]models.URL, int64, error) {
	urls, total, err := u.urlRepo.FindAllByUserId(userId, page, pageSize)
	if err != nil {
		return nil, -1, &appErrors.AppError{
			Code:    "DB_ERROR",
			Message: "Database error",
			Status:  http.StatusInternalServerError,
		}
	}
	return urls, total, nil
}

func (u *urlService) UpdateUrl(userId uint, urlId uint, req dto.UpdateUrlRequest) (*models.URL, error) {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)
	if err != nil {
		return nil, appErrors.ErrURLNotFound
	}

	if req.OriginalUrl != "" {
		url.OriginalURL = req.OriginalUrl
	}

	if req.ExpiresAt != nil {
		url.ExpiresAt = req.ExpiresAt
	}

	if req.IsActive != nil {
		url.IsActive = *req.IsActive
	}

	if err := u.urlRepo.Update(url); err != nil {
		return nil, &appErrors.AppError{
			Code:    "DB_ERROR",
			Message: "Database error",
			Status:  http.StatusInternalServerError,
		}
	}

	if err := u.cache.Delete(context.Background(), "url:"+url.ShortCode); err != nil {
		metrics.CacheInvalidationErrorsTotal.Inc()

		log.Printf(
			"failed to invalidate cache for URL %d: %v",
			urlId,
			err,
		)
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

	if err := u.cache.Delete(context.Background(), "url:"+url.ShortCode); err != nil {
		metrics.CacheInvalidationErrorsTotal.Inc()

		log.Printf(
			"failed to invalidate cache for URL %d: %v",
			urlId,
			err,
		)
	}

	return nil
}

func (u *urlService) ActivateUrl(userId uint, urlId uint) error {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)
	if err != nil {
		return appErrors.ErrURLNotFound
	}

	if err := u.urlRepo.UpdateIsActive(urlId, true); err != nil {
		return appErrors.ErrServerError
	}

	if err := u.cache.Delete(context.Background(), "url:"+url.ShortCode); err != nil {
		metrics.CacheInvalidationErrorsTotal.Inc()

		log.Printf(
			"failed to invalidate cache for URL %d: %v",
			urlId,
			err,
		)
	}

	return nil
}

func (u *urlService) DeactivateUrl(userId uint, urlId uint) error {
	url, err := u.urlRepo.FindByIdAndUserId(urlId, userId)
	if err != nil {
		return appErrors.ErrURLNotFound
	}

	if err := u.urlRepo.UpdateIsActive(urlId, false); err != nil {
		return appErrors.ErrServerError
	}

	if err := u.cache.Delete(context.Background(), "url:"+url.ShortCode); err != nil {
		metrics.CacheInvalidationErrorsTotal.Inc()

		log.Printf(
			"failed to invalidate cache for URL %d: %v",
			urlId,
			err,
		)
	}

	return nil
}

func (u *urlService) Redirect(shortCode string, ctx context.Context) (*models.URL, error) {
	key := "url:" + shortCode

	if u.cache != nil {
		cachedData, err := u.cache.Get(ctx, key)

		if err == nil {
			var cached CachedURL

			if jsonError := json.Unmarshal([]byte(cachedData), &cached); jsonError == nil {

				if cached.ExpiresAt != nil && !cached.ExpiresAt.After(time.Now()) {
					_ = u.urlRepo.UpdateIsActive(cached.ID, false)
					_ = u.cache.Delete(ctx, key)

					return nil, &appErrors.AppError{
						Code:    "URL_EXPIRES_ALREADY",
						Message: "Url already expired",
						Status:  http.StatusNoContent,
					}
				}

				return &models.URL{
					ID:          cached.ID,
					OriginalURL: cached.OriginalURL,
					ShortCode:   shortCode,
					ExpiresAt:   cached.ExpiresAt,
					IsActive:    cached.IsActive,
					IsDeleted:   cached.IsDeleted,
				}, nil
			}
		}
	}

	// Cache MISS: Query Database
	url, err := u.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return nil, appErrors.ErrURLNotFound
	}

	if url.IsDeleted || !url.IsActive {
		return nil, &appErrors.AppError{
			Code:    "URL_ALREADY_DELETED",
			Message: "Url already deleted",
			Status:  http.StatusNoContent,
		}
	}

	if url.ExpiresAt != nil && !url.ExpiresAt.After(time.Now()) {
		url.IsActive = false
		if err := u.urlRepo.UpdateIsActive(url.ID, false); err != nil {
			return nil, appErrors.ErrServerError
		}
		return nil, &appErrors.AppError{
			Code:    "URL_EXPIRES_ALREADY",
			Message: "Url already expired",
			Status:  http.StatusNoContent,
		}
	}

	cached := CachedURL{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		ExpiresAt:   url.ExpiresAt,
		IsActive:    url.IsActive,
		IsDeleted:   url.IsDeleted,
	}

	cachedData, err := json.Marshal(cached)
	if err == nil {
		if err := u.cache.Set(
			ctx,
			key,
			string(cachedData),
			10*time.Minute,
		); err != nil {
			log.Printf(
				"failed to cache URL %s: %v",
				shortCode,
				err,
			)
		}
	}

	return url, nil
}
