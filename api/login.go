package api

import (
	"encoding/json"
	"net/url"
)

type LoginResponse struct {
	Msg      string `json:"msg"`
	Success  bool   `json:"success"`
	Errors   []any  `json:"errors"`
	Warnings []any  `json:"warnings"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Client   string `json:"client"`
	Clients  []struct {
		Client string `json:"client"`
	} `json:"clients"`
	Dbaccess bool `json:"dbaccess"`
}

func (me *API) Login(username string, password string) (LoginResponse, error) {
	var loginResponse LoginResponse

	sb, err := me.PostForm(me.systemURL+"/login", url.Values{
		"forcelogin": {"1"},
		"username":   {username},
		"password":   {password},
	})
	json.Unmarshal([]byte(sb), &loginResponse)
	return loginResponse, err
}
