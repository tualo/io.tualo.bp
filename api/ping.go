package api

import "encoding/json"

type PingResponse struct {
	Msg      string `json:"msg"`
	Success  bool   `json:"success"`
	Errors   []any  `json:"errors"`
	Warnings []any  `json:"warnings"`
	Username string `json:"username"`
	Clients  []struct {
		Client string `json:"client"`
	} `json:"clients"`
	Client       string `json:"client"`
	Fullname     string `json:"fullname"`
	Gst          string `json:"gst"`
	Bkr          string `json:"bkr"`
	Gstavatar    string `json:"gstavatar"`
	Bkravatar    string `json:"bkravatar"`
	Avatar       string `json:"avatar"`
	Clientavatar string `json:"clientavatar"`
}

func (me *API) Ping() (PingResponse, error) {
	var response PingResponse
	sb, err := me.Get(me.systemURL + "dashboard/ping")
	json.Unmarshal([]byte(sb), &response)
	me.lastPing = response
	return response, err
}
