package users

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

func UsersDelete(c iris.Context) {
	id := c.Params().Get("id")
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Missing user ID"})
		return
	}

	if err := database.Db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to delete user"})
		return
	}

	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{"message": "User deleted successfully"})
}
