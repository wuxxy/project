package users

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

type CreateUserRequest struct {
	Username   string   `json:"username" validate:"required"`
	Password   string   `json:"password" validate:"required"`
	Email      string   `json:"email" validate:"required"`
	Locations  []string `json:"locations"`
	UserAgents []string `json:"user_agents"`
	AvatarURL  string   `json:"avatar_url"`
}

func UsersCreate(c iris.Context) {
	var req CreateUserRequest

	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Invalid request body"})
		return
	}

	user := models.User{
		ID:         uuid.New().String(),
		Username:   req.Username,
		Password:   req.Password,
		Email:      req.Email,
		Locations:  req.Locations,
		UserAgents: req.UserAgents,
		AvatarURL:  req.AvatarURL,
	}

	if err := database.Db.Create(&user).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to create user"})
		return
	}

	c.StatusCode(iris.StatusCreated)
	_ = c.JSON(iris.Map{
		"message": "User created successfully",
		"user":    user,
	})
}
