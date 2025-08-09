package models

type User struct {
	ID string `gorm:"primaryKey" json:"id"`
}
