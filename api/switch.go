package api

import "encoding/json"

type SwitchClientResponse struct {
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

func (me *API) Switch(toClient string) (SwitchClientResponse, error) {
	var response SwitchClientResponse
	sb, err := me.Post(me.systemURL+"dashboard/client/switch", "toclient=bwmuenchen")
	json.Unmarshal([]byte(sb), &response)
	return response, err
}
