package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/auth/database"
	"github.com/wuxxy/auth/models"
)

func Me(c iris.Context) {
	// Get user ID from context
	userID := c.Values().Get("user_id")
	if userID == nil {
		c.StatusCode(iris.StatusUnauthorized)
		_ = c.JSON(iris.Map{"error": "Unauthorized"})
		return
	}

	// Fetch user details from database
	var user models.User
	if err := database.Db.
		Omit("Password").
		Preload("Sessions").
		First(&user, "users.id = ?", userID).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to fetch user details"})
		return
	}
	// Return user details
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(user)
}
