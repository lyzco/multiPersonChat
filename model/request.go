package model

type Request struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type JoinData struct {
	UserName string `json:"userName"`
	Room     string `json:"room"`
}
