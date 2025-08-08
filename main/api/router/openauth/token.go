package openauth

import "github.com/kataras/iris/v12"

type OpenAuthTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	State        string `json:"state"`
}

func OpenAuthToken(c iris.Context) {

}
