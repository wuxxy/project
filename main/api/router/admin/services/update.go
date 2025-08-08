package services

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

func ServicesUpdate(c iris.Context) {
	id := c.Params().Get("id")
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Missing service ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ReadJSON(&updateData); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Invalid JSON payload"})
		return
	}

	// Optional: disallow updating protected fields
	delete(updateData, "id")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")
	delete(updateData, "deleted_at")

	// Update the service using GORM
	if err := database.Db.Model(&models.Service{}).
		Where("id = ?", id).
		Updates(updateData).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to update service"})
		return
	}

	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{"message": "Service updated successfully"})
}
