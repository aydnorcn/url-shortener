package dto

import "time"

type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url"`
	ExpiresAt   *time.Time
}

type UpdateUrlRequest struct {
	OriginalUrl string `json:"original_url"`
	ExpiresAt   *time.Time
}
