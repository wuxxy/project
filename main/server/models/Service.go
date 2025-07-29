package models

type Service struct {
	ID          string `json:"id" gorm:"primaryKey"`        // Primary key for the service
	Name        string `gorm:"unique;not null" json:"name"` // Service name, must be unique and not null
	Secret      string `gorm:"not null" json:"secret"`
	RedirectUrl string `gorm:"not null" json:"redirect_url"`
}
