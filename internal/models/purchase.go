package models

type Purchase struct {
	Date    string     `json:"created_at"`
	Creator SimpleUser `json:"creator"`
	Type    string     `json:"type"`
	Value   string     `json:"value"`
	Amount  float64    `json:"amount"`
}

type SimpleUser struct {
	ID          int    `json:"id"`
	Avatar      string `json:"avatar"`
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
}
