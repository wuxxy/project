package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/auth/database"
	"github.com/wuxxy/auth/models"
)

func Logout(c iris.Context) {
	// Clear the refresh token cookie
	sessionID := c.Values().Get("session_id")
	if sessionID == nil {
		c.StatusCode(iris.StatusUnauthorized)
		_ = c.JSON(iris.Map{"error": "Unauthorized"})
		return
	}
	err := database.Db.Delete(&models.Session{}, "id = ?", sessionID).Error
	if err != nil {
		_ = c.JSON(iris.Map{"error": "Failed to delete session and logout"})
		return
	}
	c.RemoveCookie("refresh")
	// Respond with a success message
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{"success": "true"})
}
