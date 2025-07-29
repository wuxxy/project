package users

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

func UsersUpdate(c iris.Context) {
	id := c.Params().Get("id")
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Missing user ID"})
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
	delete(updateData, "sessions")

	// Update the user using GORM
	if err := database.Db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(updateData).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to update user"})
		return
	}

	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{"message": "User updated successfully"})
}
