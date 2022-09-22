package model

// Response struct
type Response struct {
	Type int         `json:"type"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type HeartCheckData struct {
	Id int `json:"id"`
}

type ResponseDataOO struct {
	UserName string `json:"userName"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

type ResponseRoomInfo struct {
	RoomName string `json:"roomName"`
	RoomId   string `json:"roomId"`
	Users    string `json:"users"`
}
