package openauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

// OAuth2Authorize handles the OAuth2 authorization flow.
// This is only for server-side application since the frontend will handle the redirect.
// It is only a middleman to handle the OAuth2 flow.
type OpenAuthRequest struct {
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
	State       string `json:"state"`
}

func OpenAuthAuthorize(c iris.Context) {
	var req OpenAuthRequest
	if err := c.ReadJSON(&req); err != nil {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}

	var service models.Service
	if err := database.Db.Where("id = ?", req.ClientId).First(&service).Error; err != nil {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}
	if req.RedirectUri != service.RedirectUrl {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}
	now := time.Now()
	sessionId := c.Values().Get("session_id")
	var session models.Session
	if err := database.Db.Where("id = ?", sessionId).First(&session).Error; err != nil {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	// Code: 2 min
	codeClaims := jwt.MapClaims{
		"sub": sessionId,
		"exp": now.Add(2 * time.Minute).Unix(),
		"iat": now.Unix(),
	}
	newJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, codeClaims)
	code, err := newJWT.SignedString([]byte(service.Secret))
	c.StatusCode(iris.StatusOK)
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{
		"code":         code,
		"expires_at":   codeClaims["exp"],
		"redirect_uri": req.RedirectUri,
		"state":        req.State,
	})
}
