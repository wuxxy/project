package auth

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/wuxxy/auth/database"
	"github.com/wuxxy/auth/models"
	"github.com/wuxxy/auth/tokens"
)

// Token Handler for generating a new access token using refresh token.
// This function will verify the refresh token, generate a new access token,
// and return it in the response. If the refresh token is invalid or expired,
// it will return an error response and clear the refresh token cookie.
func Token(c iris.Context) {
	refreshToken := c.GetCookie("refresh")
	if refreshToken == "" {
		c.StatusCode(iris.StatusUnauthorized)
		_ = c.JSON(iris.Map{"error": "Refresh token is required"})
		return
	}
	sessionId, err := tokens.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.StatusCode(iris.StatusUnauthorized)
		_ = c.JSON(iris.Map{"error": "Invalid or expired refresh token"})
		c.RemoveCookie("refresh")
		return
	}
	// Verifies session ID still exists in the database
	session := models.Session{
		ID: sessionId,
	}
	if err := database.Db.First(&session).Error; err != nil {
		c.StatusCode(iris.StatusUnauthorized)
		_ = c.JSON(iris.Map{"error": "Session not found"})
		c.RemoveCookie("refresh")
		return
	}
	session.LastUsed = time.Now()
	if err := database.Db.Save(&session).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not update session"})
		return
	}
	// Check if database session is expired
	if session.ExpiresAt.Before(time.Now()) {
		_ = c.JSON(iris.Map{"error": "Session expired"})
		c.RemoveCookie("refresh")
		return
	}
	// Generate a new access token
	accessToken, err := tokens.CreateAccessToken(session.ID, session.UserID)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not generate access token"})
		return
	}
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{
		"success":      true,
		"access_token": accessToken,
	})
}
