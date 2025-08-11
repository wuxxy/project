package user

import (
	"log"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuxxy/project/echo/ipc"
)

type MeResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	AvatarURL string `json:"avatar_url"`
	Suspended bool   `json:"suspended"`
	Disable   bool   `json:"disable"`
	Premium   bool   `json:"premium"`
	IsAdmin   bool   `json:"is_admin"`
	Error     string `json:"error"`
}

func Me(c iris.Context) {
	log.Print("Fetching user details for user_id: ", c.Values().GetString("user_id"))
	msgPackResponse, err := ipc.NC.Request("auth.get_user", []byte(c.Values().GetString("user_id")), 2*time.Second)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		log.Print(err.Error())
		_ = c.JSON(iris.Map{"error": "Couldn't fetch user details"})
		return
	}
	var res MeResponse
	err = msgpack.Unmarshal(msgPackResponse.Data, &res)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		log.Print("Error unmarshalling response: ", err.Error())
		_ = c.JSON(iris.Map{"error": "Internal server error"})
		return
	}
	if res.Error != "" {
		log.Print("Error from IPC: ", res.Error)
		c.StatusCode(iris.StatusInternalServerError)
		_ = c.JSON(iris.Map{"error": res.Error})
		return
	}
	log.Print(res)
	c.StatusCode(200)
	_ = c.JSON(res)

}
