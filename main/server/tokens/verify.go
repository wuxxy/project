package tokens

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyAccessToken parses and validates an access token, returning sessionID and userID
func VerifyAccessToken(tokenStr string) (sessionId string, userId string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Enforce HMAC-SHA256 signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token claims")
	}

	// Check type
	if typ, ok := claims["type"].(string); !ok || typ != "access" {
		return "", "", errors.New("not an access token")
	}

	// Extract values
	sid, sidOk := claims["sub"].(string)
	uid, uidOk := claims["uid"].(string)

	if !sidOk || !uidOk {
		return "", "", errors.New("missing required claims")
	}

	return sid, uid, nil
}
func VerifyRefreshToken(tokenStr string) (sessionId string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Enforce HMAC-SHA256 signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	// Check type
	if typ, ok := claims["type"].(string); !ok || typ != "refresh" {
		return "", errors.New("not a refresh token")
	}

	// Extract values
	sid, sidOk := claims["sub"].(string)

	if !sidOk {
		return "", errors.New("missing required claim")
	}

	return sid, nil
}
