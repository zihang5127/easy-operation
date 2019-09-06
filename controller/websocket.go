package controller

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/zihang5127/easy-operation/module/channel"
	"time"
)

// WebSocket 控制器
type WebSocketController struct {
	BaseController
}

var upgrader = websocket.Upgrader{}

func (c *WebSocketController) Ws() {

	chann := channel.GetChannel()
	oc := channel.GetOverChannel()

	c.TplName = "index/log.html"
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)

	if err != nil {
		logs.Error("%s", err)
	}

	var buffer bytes.Buffer
	go func() {
		for {
			timeout := time.NewTimer(time.Second * time.Duration(60*30))
			select {
			case temp := <-chann:
				buffer.Write(temp)
				//c.Write(ws, values)
			case <-timeout.C:
				buffer.Write([]byte("timeout"))
				//c.Write(ws, []byte("timeout"))
			case <-oc:
				goto out
			}
		}
	out:
		fmt.Println(string(buffer.Bytes()))
		buffer.Write([]byte("\n\n\n\nThe command was executed successfully !!!"))
		c.Write(ws, buffer.Bytes())
	}()

}

func (c *WebSocketController) Write(ws *websocket.Conn, msg []byte) {
	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		_ = ws.Close()
	}
}
