package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at"`
	IP        string    `json:"ip" gorm:"not null"`
	UserAgent string    `json:"user_agent" gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LastUsed  time.Time `json:"last_used" gorm:"default:CURRENT_TIMESTAMP"`
}
