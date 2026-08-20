package dto

import "time"

type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ExpiresAt   *time.Time
	CustomAlias *string `json:"custom_alias"`
}

type UpdateUrlRequest struct {
	OriginalUrl string     `json:"original_url"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    *bool      `json:"is_active"`
}

type UrlResponse struct {
	Id          int       `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortCode   int       `json:"short_code"`
	ShortUrl    string    `json:"short_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
}
