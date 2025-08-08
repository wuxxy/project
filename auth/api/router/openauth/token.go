package openauth

import (
	"log"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
	"github.com/wuxxy/project/main/tokens"
)

type OpenAuthTokenRequest struct {
	ServiceId     string `json:"service_id"`
	ServiceSecret string `json:"service_secret"`
	State         string `json:"state"`
	Code          string `json:"code"`
}

func OpenAuthToken(c iris.Context) {
	var req OpenAuthTokenRequest
	if err := c.ReadJSON(&req); err != nil {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// Validate the service ID and secret
	var service models.Service
	if err := database.Db.Where("id = ? AND secret = ?", req.ServiceId, req.ServiceSecret).First(&service).Error; err != nil {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	// Verify the code and state
	// Get the session ID from the code
	sessionId, state, err := tokens.ParseCodeJWT(req.Code, req.ServiceSecret)

	if err != nil {
		log.Println("JWT verification failed:", err)
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}
	// Check if the state matches
	log.Println("JWT State, Request State", state, req.State)
	if state != req.State {
		c.StopWithStatus(iris.StatusForbidden)
		return
	}
	var session models.Session
	if err = database.Db.Where("id = ?", sessionId).First(&session).Error; err != nil {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}
	accessToken, refreshToken, err := tokens.CreatePair(session.ID, session.UserID)
	c.SetCookieKV("refresh", refreshToken,
		iris.CookieExpires(24*time.Hour),
		iris.CookieHTTPOnly(true),
		iris.CookiePath("/"),
		iris.CookieSameSite(http.SameSiteLaxMode),
	)
	// Send access token as success response
	c.StatusCode(200)
	_ = c.JSON(iris.Map{
		"success":      true,
		"access_token": accessToken,
	})

}
