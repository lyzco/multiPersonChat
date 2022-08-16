package main

import (
	"chatroom/client/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func help() {
	str := "\n   ___  _____     __                       \n  / __\\/__   \\   /__\\ ___   ___  _ __ ___  \n / /     / /\\/  / \\/// _ \\ / _ \\| '_ ` _ \\ \n/ /___  / /    / _  \\ (_) | (_) | | | | | |\n\\____/  \\/     \\/ \\_/\\___/ \\___/|_| |_| |_|\n                                           \n"
	fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, str, 0x1B)

}

func initWs() *websocket.Conn {
	var wsUrl = "ws://127.0.0.1:9502/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{})
	if err != nil {
		panic(err)
	}

	return ws
}

func initLog() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func main() {
	help()
	initLog()

	var userName string
	var operateType int
	var roomId int
	var roomName string
	var requestData model.Request
	var ws *websocket.Conn
	var sendMsg string
	//var consoleLog chan map[string]string

	for {
		fmt.Printf("请输入用户名：")
		fmt.Scanf("%s", &userName)
		if userName != "" {
			break
		}
	}

	fmt.Println()
optCheck:
	for {
		fmt.Printf("请选择：(1)加入房间 (2)创建房间 :")
		fmt.Scanf("%d", &operateType)
		if operateType == 1 {
			for {
				fmt.Printf("请输入房间号：")
				fmt.Scanf("%d", &roomId)
				if roomId > 9999 {
					ws = initWs()
					requestData = model.Request{
						Type: model.MsgTypeJoin, Data: map[string]interface{}{
							"userName": userName,
							"roomId":   roomId,
						}}
					break optCheck
				}
			}
		} else if operateType == 2 {
			for {
				fmt.Printf("请输入房间名称：")
				fmt.Scanf("%s", &roomName)
				if roomName != "" {
					ws = initWs()
					requestData = model.Request{
						Type: model.MsgTypeCreate,
						Data: map[string]string{
							"userName": userName,
							"roomName": roomName,
						}}
					break optCheck
				}
			}
		}
	}

	requestString, _ := json.Marshal(requestData)
	err := ws.WriteMessage(websocket.TextMessage, requestString)
	if err != nil {
		zap.S().Errorf("send ws message failed, with error %s", err.Error())
		return
	}

	go func() {
		for {
			if _, Data, err := ws.ReadMessage(); err != nil {
				zap.S().Errorf("get ws message failed, with error %s", err.Error())
				return
			} else {

				//zap.S().Infof(" messgae %s", string(Data))
				wsResponse := model.Response{}

				if err := json.Unmarshal(Data, &wsResponse); err != nil {
					zap.S().Errorf("unmarshal ws req error: %s", err.Error())
					return
				}
				switch wsResponse.Type {
				case model.MsgTypeBot:
					msgData := wsResponse.Data.(map[string]interface{})
					if _, ok := msgData["roomId"]; ok {
						msgData["message"] = "创建房间成功" + msgData["roomId"].(string)
					}
					log.Output(0, "bot:"+msgData["message"].(string))
				case model.MsgTypeBusiness:
					msgData := wsResponse.Data.(map[string]interface{})
					err := log.Output(0, msgData["userName"].(string)+":"+msgData["message"].(string))
					if err != nil {
						return
					}
				case model.MsgTypeHeart:
					err := ws.WriteMessage(websocket.TextMessage, Data)
					if err != nil {
						return
					}
				case model.MSGTypeError:
					zap.S().Error("get ws message failed")
					os.Exit(1)
				default:
				}
			}

		}
	}()

	for {
		fmt.Scanf("%s", &sendMsg)
		if sendMsg != "" {
			requestData = model.Request{
				Type: model.MsgTypeBusiness,
				Data: map[string]string{
					"message": sendMsg,
				},
			}
			requestString, _ := json.Marshal(requestData)
			err := ws.WriteMessage(websocket.TextMessage, requestString)
			if err != nil {
				return
			}
		}
	}

	//for {
	//	select {
	//	case c := <-consoleLog:
	//		//fmt.Println(c)
	//		//zap.
	//	}
	//}
}
