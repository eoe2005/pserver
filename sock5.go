package pserver

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

func Sock5(ws *websocket.Conn) {
	defer ws.Close()

	_, d, e := ws.ReadMessage()
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}
	Debug("read 1 WS %d", len(d))
	e = ws.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x00})
	if e != nil {
		Debug("write ws 失败 %s", e.Error())
		return
	}

	_, h2, e := ws.ReadMessage()
	Debug("read 2 WS %d", len(h2))
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}
	h2 = h2[0:4]
	addr := ""
	e = ws.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x00, 0x00, h2[3]})
	if e != nil {
		Debug("write ws 失败 %s", e.Error())
		return
	}
	Debug("3 -> %d", h2[3])
	switch h2[3] {
	case 0x01:
		Debug("001")
		// _, v, _ := ws.ReadMessage()
		v := d[7:11]
		addr = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
		Debug("链接地址 %s", addr)
		ws.WriteMessage(websocket.BinaryMessage, v)
	case 0x04:
		// _, v, _ := ws.ReadMessage()
		Debug("002")
		v := d[7:23]
		ws.WriteMessage(websocket.BinaryMessage, v)
	case 0x03:
		Debug("003")
		// _, l, _ := ws.ReadMessage()
		l := d[7:8]
		v := d[8 : 8+l[0]]
		ws.WriteMessage(websocket.BinaryMessage, l)
		ws.WriteMessage(websocket.BinaryMessage, v)
		addr = string(v)
		Debug("链接地址 %s", addr)
	}
	if addr == "" {
		Debug("read not addr")
		return
	}
	_, p, e := ws.ReadMessage()
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}

	port := (int(p[0]) << 8) + int(p[1])
	Debug("链接端口 %d - %d %d %s", port, int(p[0])<<8, p[1], string(p))

	desc, e := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		Debug("建立远端链接 失败 %s", e.Error())
		return
	}

	e = ws.WriteMessage(websocket.BinaryMessage, p)
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}

	defer desc.Close()
	TraceDate(desc, ws)
}
