package services

import (
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

func ServicesReadAll(c iris.Context) {
	// This function is a placeholder for the services read all endpoint.
	// It should handle the logic to retrieve all services.
	// Currently, it does not perform any operations.
	var services []models.Service
	if err := database.Db.
		Find(&services).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": "Failed to fetch user details"})
		return
	}
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(services)
}

func ServicesStruct(c iris.Context) {
	// Get the type info for the struct
	serviceType := reflect.TypeOf(models.Service{})

	structInfo := make(map[string]string)

	for i := 0; i < serviceType.NumField(); i++ {
		field := serviceType.Field(i)

		// Optional: skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Use JSON tag if available, otherwise use field name
		jsonTag := field.Tag.Get("json")
		fieldName := jsonTag
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}

		// Store type name
		structInfo[fieldName] = field.Type.String()
	}

	_ = c.JSON(structInfo)
}
