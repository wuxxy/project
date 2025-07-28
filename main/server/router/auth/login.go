package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/9ssi7/turnstile"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
	"github.com/wuxxy/project/main/tokens"
	"gorm.io/gorm"
)

// LoginRequest The structure of the request body for user login
type LoginRequest struct {
	Username       string `json:"username" validate:"required,min=4,max=14"`
	Password       string `json:"password" validate:"required"`
	RecaptchaToken string `json:"turnstile" validate:"required"`
}

// Login Handler for user login
// Validates the request, checks the CAPTCHA, verifies the user credentials
// and creates a session for the user.
// Returns access token and sets refresh token as a cookie.
func Login(c iris.Context) {
	// Assign request struct to body
	var req LoginRequest

	// Parse input as request body
	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Invalid request body"})
		return
	}
	var err error
	if c.GetHeader("SiteTestingKey") != os.Getenv("STK") {
		// Verify CAPTCHA to prevent spam
		srv := turnstile.New(turnstile.Config{
			Secret: os.Getenv("TURNSTILE_SECRET"),
		})
		ok, err := srv.Verify(c, req.RecaptchaToken, c.Host())
		if err != nil || !ok {
			c.StatusCode(iris.StatusInternalServerError)
			_ = c.JSON(iris.Map{"error": "Failed to verify CAPTCHA"})
			return
		}
	}
	// Make sure the request body is valid and meets requirements
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errs := make(map[string]string)
		for _, e := range validationErrors {
			field := e.Field()
			switch field {
			case "Username":
				errs["username"] = "Username must be 4–14 characters"
			case "Password":
				errs["password"] = "Password must exist"
			case "RecaptchaToken":
				errs["turnstile"] = "Missing CAPTCHA token"
			}
		}
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"errors": errs})
		return
	}
	// Check if user exists (get user)
	var user models.User
	err = database.Db.
		Where("username = ?", req.Username).
		First(&user).Error
	// User doesn't exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = c.JSON(iris.Map{"error": "Invalid username or password"})
		return
	}
	// Failed to fetch user
	if err != nil {
		_ = c.JSON(iris.Map{"error": "Database error", "details": err.Error()})
		return
	}
	// Make sure the password is also valid
	passwordValid := ValidatePassword(req.Password)
	if passwordValid != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": passwordValid.Error()})
		return
	}

	isCorrectPassword, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not verify password"})
		return
	}
	if !isCorrectPassword {
		c.StatusCode(iris.StatusUnauthorized)
		_ = c.JSON(iris.Map{"error": "Invalid username or password"})
		return
	}

	// Now you can use newUser.ID
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.RemoteAddr()
	}
	session := models.Session{
		ID:        uuid.NewString(), // if you're using UUIDs
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IP:        ip,
		UserAgent: c.GetHeader("User-Agent"),
	}
	if err := database.Db.Create(&session).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not create session"})
		return
	}

	accessToken, refreshToken, err := tokens.CreatePair(session.ID, user.ID) // if ID is UUID or convert to string

	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not generate tokens"})
		return
	}
	// Set the refresh token as a cookie
	c.SetCookieKV("refresh", refreshToken,
		iris.CookieExpires(30*24*time.Hour),
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
