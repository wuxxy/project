package services

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

func ServicesDelete(c iris.Context) {
	id := c.Params().Get("id")
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Missing service ID"})
		return
	}

	// Delete from database
	if err := database.Db.Delete(&models.Service{}, "id = ?", id).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to delete service"})
		return
	}

	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{"message": "Service deleted successfully"})
}
