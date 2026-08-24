package repository

import (
	"context"
	"url-shortener/models"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	CreateClick(ctx context.Context, click *models.URLClick) error
	GetTotalClicks(ctx context.Context, urlID uint) (int64, error)
	GetUniqueVisitors(ctx context.Context, urlID uint) (int64, error)
	GetDeviceStats(ctx context.Context, urlID uint) (map[string]int64, error)
	GetCountryStats(ctx context.Context, urlID uint) (map[string]int64, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) CreateClick(ctx context.Context, click *models.URLClick) error {
	return r.db.WithContext(ctx).Create(click).Error
}

func (r *analyticsRepository) GetTotalClicks(ctx context.Context, urlID uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.URLClick{}).
		Where("url_id = ?", urlID).
		Count(&total).Error
	return total, err
}

func (r *analyticsRepository) GetUniqueVisitors(ctx context.Context, urlID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.URLClick{}).
		Where("url_id = ? AND ip_address != ''", urlID).
		Distinct("ip_address").
		Count(&count).Error
	return count, err
}

func (r *analyticsRepository) GetDeviceStats(ctx context.Context, urlID uint) (map[string]int64, error) {
	type DeviceResult struct {
		Device string `gorm:"column:device"`
		Count  int64  `gorm:"column:count"`
	}

	var results []DeviceResult
	err := r.db.WithContext(ctx).
		Model(&models.URLClick{}).
		Select("device, count(*) as count").
		Where("url_id = ?", urlID).
		Group("device").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, res := range results {
		key := res.Device
		if key == "" {
			key = "unknown"
		}
		stats[key] = res.Count
	}

	return stats, nil
}

func (r *analyticsRepository) GetCountryStats(ctx context.Context, urlID uint) (map[string]int64, error) {
	type CountryResult struct {
		Country string `gorm:"column:country"`
		Count   int64  `gorm:"column:count"`
	}

	var results []CountryResult
	err := r.db.WithContext(ctx).
		Model(&models.URLClick{}).
		Select("country, count(*) as count").
		Where("url_id = ?", urlID).
		Group("country").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, res := range results {
		key := res.Country
		if key == "" {
			key = "unknown"
		}
		stats[key] = res.Count
	}

	return stats, nil
}
