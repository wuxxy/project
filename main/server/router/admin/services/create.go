package services

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

type CreateServiceRequest struct {
	Name        string `json:"name" validate:"required"`
	Secret      string `json:"secret" validate:"required"`
	RedirectUri string `json:"redirect_uri" validate:"required"`
}

func generateRandomSecret() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ServicesCreate(c iris.Context) {
	var req CreateServiceRequest

	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Invalid request body"})
		return
	}

	// Auto-generate secret if explicitly requested
	secret := req.Secret
	if secret == "auto()" {
		gen, err := generateRandomSecret()
		if err != nil {
			c.StatusCode(iris.StatusInternalServerError)
			_ = c.JSON(iris.Map{"error": "Failed to generate secret"})
			return
		}
		secret = gen
	}

	service := models.Service{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Secret:      secret,
		RedirectUrl: req.RedirectUri,
	}

	if err := database.Db.Create(&service).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to create service"})
		return
	}

	c.StatusCode(iris.StatusCreated)
	_ = c.JSON(iris.Map{
		"message": "Service created successfully",
		"service": service,
	})
}
