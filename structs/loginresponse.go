package structs

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
