package ipc

import (
	"github.com/wuxxy/project/main/ipc/handlers"
)

func InitHandler() {
	_, _ = NC.Subscribe("auth.verify_token", handlers.VerifyToken)
	_, _ = NC.Subscribe("auth.get_user", handlers.GetUser)
}
