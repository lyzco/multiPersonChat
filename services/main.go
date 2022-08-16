package main

import (
	"chatroom/services/ws"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	go ws.NewClientManager().Start()

	initRouter()
}

func initRouter() {
	ApiRouter := gin.Default()
	gin.SetMode("debug")
	//ApiRouter.GET("/", func(c *gin.Context) {
	//
	//	c.HTML(http.StatusOK, "index.html", struct{}{})
	//})
	//
	//ApiRouter.GET("/chat", func(c *gin.Context) {
	//	data := gin.H{
	//		"title": "chatgo",
	//	}
	//	c.HTML(http.StatusOK, "chat.html", data)
	//})

	ApiRouter.GET("/ws", ws.HandleWs)

	err := ApiRouter.Run(":" + strconv.Itoa(9502))
	if err != nil {
		panic("gin listen" + err.Error())
	}
}
