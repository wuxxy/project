package models

type User struct {
	ID            string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username      string  `json:"username" gorm:"unique;not null"`
	Stories       []Story `json:"stories" gorm:"foreignKey:user_stories;references:ID"`
	ForeignUserID string  `json:"foreign_user_id" gorm:"unique;not null"`
}
