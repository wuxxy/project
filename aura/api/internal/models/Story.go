package models

type Story struct {
	ID   string `json:"id" gorm:"primaryKey"` // Primary key for the service
	User User   `gorm:"foreignKey:user_stories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
