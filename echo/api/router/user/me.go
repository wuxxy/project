package user

import "github.com/kataras/iris/v12"

func Me(c iris.Context) {
	// Get the user ID from the context
	c.JSON("test")
}
