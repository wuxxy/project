package middleware

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/tokens"
)

func AuthMiddleware(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		_ = ctx.JSON(iris.Map{"error": "Missing Authorization header"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		ctx.StatusCode(iris.StatusUnauthorized)
		_ = ctx.JSON(iris.Map{"error": "Invalid Authorization header format"})
		return
	}
	sessionID, userID, err := tokens.VerifyAccessToken(parts[1])
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_ = ctx.JSON(iris.Map{"error": "Invalid or expired token"})
		return
	}

	// Store userID in context
	ctx.Values().Set("user_id", userID)
	ctx.Values().Set("session_id", sessionID)
	ctx.Next()
}
