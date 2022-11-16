package pserver

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

func Sock5Raw(ws *websocket.Conn) {
	Debug("开始建立链接")
	_, r, _ := ws.NextReader()
	w, _ := ws.NextWriter(websocket.BinaryMessage)
	rb := make([]byte, 3)
	r.Read(rb)
	w.Write([]byte{0x05, 0x00})
	rb = make([]byte, 4)
	r.Read(rb)
	addr := ""
	w.Write([]byte{0x05, 0x00, 0x00, rb[3]})
	Debug("具体模式 %v", rb)
	switch rb[3] {
	case 0x01:
		Debug("001 ")
		v := make([]byte, 4)
		r.Read(v)
		addr = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
		Debug("链接地址 %s", addr)
		w.Write(v)
	case 0x04:
		// _, v, _ := ws.ReadMessage()
		Debug("002")
		v := make([]byte, 16)
		r.Read(v)
		w.Write(v)
	case 0x03:
		l := make([]byte, 1)
		r.Read(l)
		v := make([]byte, l[0])
		r.Read(v)
		w.Write(l)
		w.Write(v)
		addr = string(v)
		Debug("链接地址 %s", addr)
	}
	if addr == "" {
		Debug("read not addr")
		return
	}
	p := make([]byte, 2)
	r.Read(p)
	port := (int(p[0]) << 8) + int(p[1])
	Debug("链接端口 %d - %d %d %s", port, int(p[0])<<8, p[1], string(p))

	desc, e := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		Debug("建立远端链接 失败 %s", e.Error())
		return
	}
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}

	TraceRaw(desc, ws)
}
func Sock5(ws *websocket.Conn) {

	_, d, e := ws.ReadMessage()
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}
	Debug("read 1 WS %d -> %v", len(d), d[0:30])
	e = ws.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x00})
	if e != nil {
		Debug("write ws 失败 %s", e.Error())
		return
	}
	Debug("write 2")
	_, h21, e := ws.ReadMessage()
	Debug("read 2 WS %d", len(h21))
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}
	h2 := h21[0:4]
	addr := ""
	rm := []byte{0x05, 0x00, 0x00}
	// e = ws.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x00, 0x00, h2[3]})
	// if e != nil {
	// 	Debug("write ws 失败 %s", e.Error())
	// 	return
	// }
	var p []byte
	Debug("3 -> %d -> %s", h2[3], string(h21))
	switch h2[3] {
	case 0x01:
		Debug("001 ")
		// _, v, _ := ws.ReadMessage()
		v := h21[4:8]
		addr = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
		Debug("链接地址 %s", addr)
		p = h21[8:10]
		ws.WriteMessage(websocket.BinaryMessage, append(rm, h21[3:8]...))
	case 0x04:
		// _, v, _ := ws.ReadMessage()
		Debug("002")
		ws.WriteMessage(websocket.BinaryMessage, append(rm, h21[3:20]...))
	case 0x03:
		l := h21[4:5]
		lenh := 5 + l[0]
		v := h21[5:lenh]

		p = h21[lenh : lenh+2]

		Debug("链接地址 %d -> %s %v", l[0], string(v), p)
		ws.WriteMessage(websocket.BinaryMessage, append(rm, h21[3:lenh+3]...))
		addr = string(v)
		Debug("链接地址 %s", addr)
	}
	if addr == "" {
		Debug("read not addr")
		return
	}
	// _, p, _ = ws.ReadMessage()
	port := (int(p[0]) << 8) + int(p[1])
	Debug("链接端口 %d - %d %d %s", port, int(p[0])<<8, p[1], string(p))

	desc, e := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		Debug("建立远端链接 失败 %s", e.Error())
		return
	}
	if e != nil {
		Debug("read ws 失败 %s", e.Error())
		return
	}

	TraceDate(desc, ws)
}
