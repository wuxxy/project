package models

type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
