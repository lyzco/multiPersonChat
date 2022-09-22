package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/lyzco/multiPersonChat/model"
	"go.uber.org/zap"
	"runtime/debug"
	"strconv"
	"time"
)

type Client struct {
	Addr         string
	Socket       *websocket.Conn
	SendChain    chan []byte
	RoomId       int
	Username     string
	PingId       int
	IsReposePing bool
}

func NewClient(addr string, socket *websocket.Conn) *Client {
	return &Client{
		Addr:      addr,
		Socket:    socket,
		SendChain: make(chan []byte),
	}
}

// unique key room + username
func (c *Client) GetConnId() string {
	return strconv.Itoa(c.RoomId) + c.Username
}

func (c *Client) Read() {

	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("Read msg panic: %s", string(debug.Stack()))
		}
	}()

	defer func() {
		zap.S().Infof("close clien sendChain, %v", c)
		close(c.SendChain)
	}()

	for {
		if messageType, Data, err := c.Socket.ReadMessage(); err != nil {
			zap.S().Errorf("get ws message failed, with error %s", err.Error())
			return
		} else {

			zap.S().Infof("addr %s, messageType %d, messgae %s", c.Addr, messageType, string(Data))

			HandleMsg(c, Data)
		}

	}
}

func HandleMsg(c *Client, msg []byte) {
	wsRequest := model.Request{}

	if err := json.Unmarshal(msg, &wsRequest); err != nil {
		zap.S().Errorf("unmarshal ws req error: %s", err.Error())
		return
	}

	if c.RoomId == 0 && wsRequest.Type != 30 && wsRequest.Type != 40 {
		res := model.Response{
			Type: model.MsgTypeError,
			Msg:  "请先加入房间",
		}
		resJson, _ := json.Marshal(res)
		c.SendMsg(resJson)
		Manager.CloseSocket(c)
		return
	}

	switch wsRequest.Type {
	case model.MsgTypeBusiness:
		requestData := wsRequest.Data.(map[string]interface{})
		if requestData != nil {
			message := requestData["message"].(string)
			room := Manager.GetRoomById(c.RoomId)
			room.Send(model.MsgTypeBusiness, message, c.Username, c)
		}
		break
	case model.MsgTypeHeart:
		//todo optimize
		c.IsReposePing = true
		break
	case model.MsgTypeJoin:

		requestData := wsRequest.Data.(map[string]interface{})

		if requestData != nil {
			c.Username = requestData["userName"].(string)

			roomId, _ := strconv.Atoi(strconv.FormatFloat(requestData["roomId"].(float64), 'f', model.MsgTypeError, 64))
			room := Manager.GetRoomById(roomId)

			if room == nil {
				res := model.Response{
					Type: model.MsgTypeError,
					Msg:  "房间不存在，请重试",
				}
				resJson, _ := json.Marshal(res)

				c.SendMsg(resJson)
				Manager.CloseSocket(c)
				return
			}
			_, ok := room.Users[c.Username]

			if ok {
				res := model.Response{
					Type: model.MsgTypeError,
					Msg:  "用户名重复，请重试",
				}
				resJson, _ := json.Marshal(res)
				c.SendMsg(resJson)
				Manager.CloseSocket(c)
				return
			}

			c.RoomId = roomId
			room.AddUser(c.Username, c)
			room.Send(model.MsgTypeBot, "join", c.Username, c)

		}
		go c.HeartCheck()
		break
	case model.MsgTypeCreate:
		requestData := wsRequest.Data.(map[string]interface{})
		if requestData != nil {
			roomName := requestData["roomName"].(string)
			c.Username = requestData["userName"].(string)
			room := NewRoom(roomName)
			c.RoomId = room.Id
			room.AddUser(c.Username, c)
			room.SendRoomInfo(10, c)
			Manager.AddRoom(room.Id, room)
			go c.HeartCheck()
		}
		break
	default:
		res := model.Response{
			Type: model.MsgTypeError,
			Msg:  "消息错误，请重试",
		}
		resJson, _ := json.Marshal(res)

		c.SendMsg(resJson)
		Manager.CloseSocket(c)
		break
	}
}

func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("Write msg panic: %s", string(debug.Stack()))

		}
	}()

	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		select {
		case msg, ok := <-c.SendChain:
			if !ok {
				zap.S().Error("write msg error")
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				zap.S().Error("write msg error" + err.Error())
			}
		}
	}
}

func (c *Client) SendMsg(msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("Send msg panic: %s", string(debug.Stack()))
		}
	}()
	c.SendChain <- msg
}

// HeartCheck
func (c *Client) HeartCheck() {
	data := model.Response{
		Type: 20,
		Data: nil,
	}
	for {
		c.PingId++
		if c.PingId != 1 && c.IsReposePing == false {
			c.Socket.Close()
			return
		}
		data.Data = model.HeartCheckData{Id: c.PingId}
		res, _ := json.Marshal(data)

		c.SendMsg(res)
		c.IsReposePing = false
		time.Sleep(10 * time.Second)
	}
}
