package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SQL Database model for User using GORM
type User struct {
	ID         string `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Username   string         `json:"username" gorm:"unique;not null"`
	Password   string         `json:"password" gorm:"not null"`
	Email      string         `json:"email" gorm:"unique;not null"`
	Sessions   []Session      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sessions"`
	Locations  []string       `json:"locations" gorm:"type:text[]"`
	UserAgents []string       `json:"user_agents" gorm:"type:text[]"`
	Verified   bool           `json:"verified" gorm:"default:false"`
	AvatarURL  string         `json:"avatar_url" gorm:"default:''"`
	Suspended  bool           `json:"suspended" gorm:"default:false"`
	Disable    bool           `json:"disable" gorm:"default:false"`
	Premium    bool           `json:"premium" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}
