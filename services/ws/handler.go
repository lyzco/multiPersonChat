package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWs(c *gin.Context) {
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)

	client := NewClient(conn.RemoteAddr().String(), conn)

	go client.Read()

	go client.Write()

}
