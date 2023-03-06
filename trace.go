package pserver

import (
	"io"
	"net"

	"github.com/gorilla/websocket"
)

func TraceDate(cc net.Conn, wws *websocket.Conn) {

	go func() {
		for {
			b := make([]byte, 1024*100)
			rl, e := cc.Read(b)
			if e == io.EOF {
				return
			}

			if e != nil {
				Debug("trace read:con失败 %s", e.Error())
				return
			}
			b = b[0:rl]
			Debug("sock -> ws %d -> %s", len(b), string(b))
			e = wws.WriteMessage(websocket.BinaryMessage, b)
			if e != nil {
				Debug("trace write:ws失败 %s", e.Error())
				return
			}

		}
	}()

	for {
		_, d, e := wws.ReadMessage()

		if e != nil {
			Debug("trace read :ws失败 %s", e.Error())
			return
		}
		_, e = cc.Write(d)
		Debug("ws -> sock (%v) %d->%s", e, len(d), string(d))
		if e != nil {
			Debug("trace write:con失败 %s", e.Error())
			return
		}

	}

}
