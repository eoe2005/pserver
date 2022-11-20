package pserver

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

func Http(ws *websocket.Conn) {

	_, b, e := ws.ReadMessage()
	if e != nil {
		return
	}

	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b, '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" { //https访问
		address = hostPortURL.Scheme + ":443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	fmt.Printf("read -> %v\nremote -> %s\n", string(b), string(address))
	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Printf("链接请求")
		ws.WriteMessage(websocket.BinaryMessage, []byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	} else {
		fmt.Printf("链接问题")
		server.Write(b)
	}

	TraceDate(server, ws)
}
