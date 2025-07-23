package auth

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/9ssi7/turnstile"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/auth/database"
	"github.com/wuxxy/auth/models"
	"github.com/wuxxy/auth/tokens"
)

var validate = validator.New()

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Contains at least one letter
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)

	// Contains at least one symbol (non-alphanumeric)
	hasSymbol := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	if !hasLetter || !hasSymbol {
		return errors.New("password must contain at least one letter and one symbol")
	}

	return nil
}

// RegisterRequest The structure of the request body for user login
type RegisterRequest struct {
	Email          string `json:"email" validate:"required,email"`
	Username       string `json:"username" validate:"required,min=4,max=14"`
	Password       string `json:"password" validate:"required"`
	RecaptchaToken string `json:"turnstile" validate:"required"`
}

// Register Handler for user registration
// Validates the request, checks the CAPTCHA, verifies if the user already exists,
// hashes the password, creates a new user and session, and creates session with tokens
// Returns access token and sets refresh token as a cookie.
func Register(c iris.Context) {
	// Assign request struct to body
	var req RegisterRequest

	// Parse input as request body
	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Invalid request body"})
		return
	}
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
	// Make sure the request body is valid and meets requirements
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errs := make(map[string]string)
		for _, e := range validationErrors {
			field := e.Field()
			switch field {
			case "Email":
				errs["email"] = "Invalid email format"
			case "Username":
				errs["username"] = "Username must be 4–14 characters"
			case "Password":
				errs["password"] = "Password must exist and be at least 8 characters with at least one letter and one symbol"
			case "RecaptchaToken":
				errs["turnstile"] = "Missing CAPTCHA token"
			}
		}
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"errors": errs})
		return
	}
	// Check if user with email or username already exists
	var existingUser models.User
	err = database.Db.
		Where("email = ? OR username = ?", req.Email, req.Username).
		First(&existingUser).Error

	if err == nil {
		// A user exists
		_ = c.JSON(iris.Map{"error": "Invalid or already used credentials"})
		return
	}
	// Make sure the password is also valid
	passwordValid := ValidatePassword(req.Password)
	if passwordValid != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": passwordValid.Error()})
		return
	}

	params := &argon2id.Params{
		Memory:      256 * 1024,              // 256 MB RAM usage
		Iterations:  3,                       // 3 passes over memory
		Parallelism: uint8(runtime.NumCPU()), // Leverage available cores
		SaltLength:  16,                      // Standard secure salt
		KeyLength:   32,                      // 256-bit hash
	}

	passwordHash, err := argon2id.CreateHash(req.Password, params)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not process password"})
		return
	}
	newUser := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: passwordHash,
	}

	result := database.Db.Create(&newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate") {
			c.StatusCode(iris.StatusConflict)
			_ = c.JSON(iris.Map{"error": "Invalid or already used credentials"})
			return
		}
	}

	if result.Error != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not create user"})
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
		UserID:    newUser.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IP:        ip,
		UserAgent: c.GetHeader("User-Agent"),
	}
	if err := database.Db.Create(&session).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Could not create session"})
		return
	}

	accessToken, refreshToken, err := tokens.CreatePair(session.ID, newUser.ID) // if ID is UUID or convert to string

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
