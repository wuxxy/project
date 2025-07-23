package tokens

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(strings.TrimSpace(os.Getenv("JWTSECRET")))

func CreatePair(sessionId string, userId string) (accessToken string, refreshToken string, err error) {
	now := time.Now()

	// Access Token: 15 min
	accessClaims := jwt.MapClaims{
		"sub":  sessionId,
		"uid":  userId,
		"exp":  now.Add(15 * time.Minute).Unix(),
		"iat":  now.Unix(),
		"type": "access",
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token: 7 days
	refreshClaims := jwt.MapClaims{
		"sub":  sessionId,
		"exp":  now.Add(7 * 24 * time.Hour).Unix(),
		"iat":  now.Unix(),
		"type": "refresh",
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJWT.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func CreateAccessToken(sessionId string, userId string) (accessToken string, err error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"sub":  sessionId,
		"uid":  userId,
		"exp":  now.Add(15 * time.Minute).Unix(),
		"iat":  now.Unix(),
		"type": "access",
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
