package pserver

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

func TraceDate(c net.Conn, ws *websocket.Conn) {
	wg := &sync.WaitGroup{}
	Debug("trace 开始转发")
	wg.Add(2)
	go func(c net.Conn, ws *websocket.Conn, wg *sync.WaitGroup) {
		for {
			_, d, e := ws.ReadMessage()
			Debug("readWs %s", string(d))
			if e != nil {
				Debug("trace read :ws失败 %s", e.Error())
				wg.Done()
				return
			}
			_, e = c.Write(d)
			if e != nil {
				Debug("trace write:con失败 %s", e.Error())
				wg.Done()
				return
			}
		}
	}(c, ws, wg)
	go func(c net.Conn, ws *websocket.Conn, wg *sync.WaitGroup) {
		for {
			b := make([]byte, 1024*10)
			_, e := c.Read(b)
			if e != nil {
				Debug("trace read:con失败 %s", e.Error())
				wg.Done()
				return
			}
			Debug("readCon %s", string(b))
			if ws.WriteMessage(websocket.BinaryMessage, b) != nil {
				Debug("trace write:ws失败 %s", e.Error())
				wg.Done()
				return
			}
		}
	}(c, ws, wg)
	wg.Wait()
}
