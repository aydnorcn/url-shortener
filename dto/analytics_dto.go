package dto

import "time"

type ClickEvent struct {
	URLID     uint      `json:"url_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Country   string    `json:"country"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
}

type AnalyticsResponse struct {
	TotalClicks    int64            `json:"total_clicks"`
	UniqueVisitors int64            `json:"unique_visitors"`
	Devices        map[string]int64 `json:"devices"`
	Countries      map[string]int64 `json:"countries"`
}
