package handlers

import (
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

type UserDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	AvatarURL string `json:"avatar_url"`
	Suspended bool   `json:"suspended"`
	Disable   bool   `json:"disable"`
	Premium   bool   `json:"premium"`
	IsAdmin   bool   `json:"is_admin"`
	Error     string `json:"error,omitempty"`
}

func GetUser(msg *nats.Msg) {
	var resp UserDTO

	if len(msg.Data) == 0 {
		resp.Error = "Missing ID"
		b, _ := msgpack.Marshal(&resp)
		_ = msg.Respond(b)
		return
	}

	var u models.User
	if err := database.Db.
		// Only select the fields you need; avoid gorm.DeletedAt traveling
		Select("id", "created_at", "updated_at", "deleted_at", "username",
			"email", "verified", "avatar_url", "suspended", "disable", "premium", "is_admin").
		Where("users.id = ?", string(msg.Data)).
		First(&u).Error; err != nil {

		resp.Error = "Couldn't find user"
		b, _ := msgpack.Marshal(&resp)
		_ = msg.Respond(b)
		return
	}

	resp = UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Verified:  u.Verified,
		AvatarURL: u.AvatarURL,
		Suspended: u.Suspended,
		Disable:   u.Disable,
		Premium:   u.Premium,
		IsAdmin:   u.IsAdmin,
	}

	b, _ := msgpack.Marshal(&resp)
	_ = msg.Respond(b)
}
