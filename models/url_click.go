package models

import "time"

type URLClick struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URLID     uint      `gorm:"not null;index" json:"url_id"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	Referer   string    `gorm:"type:text" json:"referer"`
	Country   string    `gorm:"type:varchar(100)" json:"country"`
	Device    string    `gorm:"type:varchar(50)" json:"device"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (URLClick) TableName() string {
	return "url_clicks"
}
