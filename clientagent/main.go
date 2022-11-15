package main

import (
	"net"
	"net/url"

	"github.com/eoe2005/pserver"
	"github.com/gorilla/websocket"
)

func main() {
	l, e := net.Listen("tcp", ":8887")
	pserver.Debug("开始监听服务")
	if e != nil {
		pserver.Debug("服务启动异常 %s", e.Error())
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			pserver.Debug("建立链接异常 %s", e.Error())
			continue
		}
		ws := GetWs()
		pserver.TraceDate(c, ws)
	}
}
func GetWs() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8888", Path: "/sock5"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		pserver.Debug("建立ws失败 %s", err.Error())
		return nil
	}
	return c
}
