package ws

import (
	"chatroom/services/model"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Room struct {
	UserLock sync.RWMutex
	Name     string
	Id       int
	Users    map[string]*Client
}

var id = 10000

func NewRoom(name string) *Room {
	id++
	return &Room{
		Name:  name,
		Users: map[string]*Client{},
		Id:    id,
	}
}

func (r *Room) AddUser(username string, c *Client) {
	r.UserLock.Lock()
	defer r.UserLock.Unlock()
	r.Users[username] = c
}

func (r *Room) DeleteUser(username string) {
	r.UserLock.Lock()
	defer r.UserLock.Unlock()
	if _, ok := r.Users[username]; ok {
		delete(r.Users, username)
	}
}

func (r *Room) GetUserList() []string {
	var ret []string

	for _, c := range r.Users {
		ret = append(ret, c.Username)
	}
	return ret
}

func (r *Room) Send(messageType int, message, username string, c *Client) {
	for _, v := range r.Users {
		if v != c {
			resp := model.Response{
				Type: messageType,
				Msg:  "success",
				Data: model.ResponseDataOO{
					UserName: username,
					Message:  message,
					Time:     time.Now().Format("2006-01-02 15:04:05"),
				},
			}
			respByte, err := json.Marshal(resp)
			if err != nil {
				zap.S().Errorf("marshal send msg at room send: %s", err.Error())
				return
			}
			v.SendChain <- respByte
		}
	}
}

func (r *Room) SendRoomInfo(msgType int, c *Client) {
	users := r.GetUserList()
	for _, v := range r.Users {

		resp := model.Response{
			Type: model.MsgTypeBot,
			Msg:  "success",
			Data: model.ResponseRoomInfo{
				RoomId:   strconv.Itoa(r.Id),
				RoomName: r.Name,
				Users:    strings.Join(users, ","),
			},
		}

		respByte, err := json.Marshal(resp)
		if err != nil {
			zap.S().Errorf("marshal send msg at room info send: %s", err.Error())
			return
		}
		v.SendChain <- respByte

	}
}
