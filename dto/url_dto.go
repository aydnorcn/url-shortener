package dto

import "time"

type CreateUrlRequest struct {
	OriginalUrl string     `json:"original_url" validate:"required,url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CustomAlias *string    `json:"custom_alias,omitempty" validate:"omitempty"`
}

type UpdateUrlRequest struct {
	OriginalUrl string     `json:"original_url,omitempty" validate:"omitempty,url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

type UrlResponse struct {
	Id          int       `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortCode   int       `json:"short_code"`
	ShortUrl    string    `json:"short_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
}
