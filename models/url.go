package models

import "time"

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	OwnerUser   User   `gorm:"not null; foreignKey:OwnerUserId"`
	OriginalURL string `gorm:"not null"`
	ShortCode   string `gorm:"uniqueIndex;not null"`
	ExpiresAt   *time.Time
	IsActive    bool `gorm:"default:true"`
	IsDeleted   bool `gorm:"default:false;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
