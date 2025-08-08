package users

import (
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

func UsersReadAll(c iris.Context) {
	var users []models.User
	if err := database.Db.
		Preload("Sessions").
		Find(&users).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to fetch user details"})
		return
	}
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(users)
}

func UsersStruct(c iris.Context) {
	userType := reflect.TypeOf(models.User{})

	structInfo := make(map[string]string)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		fieldName := jsonTag
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}

		structInfo[fieldName] = field.Type.String()
	}

	_ = c.JSON(structInfo)
}
