package models

type User struct {
	ID      string `gorm:"primaryKey" json:"id"`
	Friends []User `gorm:"many2many:user_friends;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"friends,omitempty"`
}
