package service

import (
	"context"
	"time"
	"url-shortener/appErrors"
	"url-shortener/dto"
	"url-shortener/models"
	"url-shortener/repository"
)

type AnalyticsService interface {
	RecordClick(ctx context.Context, event dto.ClickEvent) error
	GetAnalytics(ctx context.Context, userID uint, urlID uint) (*dto.AnalyticsResponse, error)
	GetTotalClicks(ctx context.Context, urlID uint) (int64, error)
	GetUniqueVisitors(ctx context.Context, urlID uint) (int64, error)
	GetDeviceStats(ctx context.Context, urlID uint) (map[string]int64, error)
	GetCountryStats(ctx context.Context, urlID uint) (map[string]int64, error)
}

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	urlRepo       repository.UrlRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, urlRepo repository.UrlRepository) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
		urlRepo:       urlRepo,
	}
}

func (s *analyticsService) RecordClick(ctx context.Context, event dto.ClickEvent) error {
	click := &models.URLClick{
		URLID:     event.URLID,
		IPAddress: event.IPAddress,
		UserAgent: event.UserAgent,
		Referer:   event.Referer,
		Country:   event.Country,
		Device:    event.Device,
		CreatedAt: event.Timestamp,
	}
	if click.CreatedAt.IsZero() {
		click.CreatedAt = time.Now()
	}
	return s.analyticsRepo.CreateClick(ctx, click)
}

func (s *analyticsService) GetAnalytics(ctx context.Context, userID uint, urlID uint) (*dto.AnalyticsResponse, error) {
	// Verify that the URL exists and is owned by the requesting user
	_, err := s.urlRepo.FindByIdAndUserId(urlID, userID)
	if err != nil {
		return nil, appErrors.ErrURLNotFound
	}

	totalClicks, err := s.GetTotalClicks(ctx, urlID)
	if err != nil {
		return nil, appErrors.ErrServerError
	}

	uniqueVisitors, err := s.GetUniqueVisitors(ctx, urlID)
	if err != nil {
		return nil, appErrors.ErrServerError
	}

	deviceStats, err := s.GetDeviceStats(ctx, urlID)
	if err != nil {
		return nil, appErrors.ErrServerError
	}
	if deviceStats == nil {
		deviceStats = make(map[string]int64)
	}

	countryStats, err := s.GetCountryStats(ctx, urlID)
	if err != nil {
		return nil, appErrors.ErrServerError
	}
	if countryStats == nil {
		countryStats = make(map[string]int64)
	}

	return &dto.AnalyticsResponse{
		TotalClicks:    totalClicks,
		UniqueVisitors: uniqueVisitors,
		Devices:        deviceStats,
		Countries:      countryStats,
	}, nil
}

func (s *analyticsService) GetTotalClicks(ctx context.Context, urlID uint) (int64, error) {
	return s.analyticsRepo.GetTotalClicks(ctx, urlID)
}

func (s *analyticsService) GetUniqueVisitors(ctx context.Context, urlID uint) (int64, error) {
	return s.analyticsRepo.GetUniqueVisitors(ctx, urlID)
}

func (s *analyticsService) GetDeviceStats(ctx context.Context, urlID uint) (map[string]int64, error) {
	return s.analyticsRepo.GetDeviceStats(ctx, urlID)
}

func (s *analyticsService) GetCountryStats(ctx context.Context, urlID uint) (map[string]int64, error) {
	return s.analyticsRepo.GetCountryStats(ctx, urlID)
}
