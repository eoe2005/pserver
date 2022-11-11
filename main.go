package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8888"
	}

	http.HandleFunc("/p", func(w http.ResponseWriter, r *http.Request) {
		log.Println("服务器的链接来了")
		testa(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("这是一个私密频道，你不要访问"))
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	log.Println("程序启动")
	http.ListenAndServe(":"+p, nil)

	// l, e := net.Listen("tcp", ":"+p)
	// if e != nil {
	// 	panic("监听端口失败")
	// }
	// for {
	// 	con, e := l.Accept()
	// 	if e != nil {
	// 		fmt.Printf("链接数据失败 %s\n", e.Error())
	// 	} else {
	// 		go handleClientRequest(con)
	// 	}
	// }
}

func testa(w http.ResponseWriter, r *http.Request) {
	con, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Printf("链接失败 %s\n", err.Error())
		return
	} else {
		log.Println("生成链接成功")
	}
	challengeKey := r.Header.Get("Sec-WebSocket-Key")
	h := sha1.New()
	h.Write([]byte(challengeKey))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	okey := base64.StdEncoding.EncodeToString(h.Sum(nil))
	send := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\nSec-WebSocket-Protocol: chat\r\nSec-WebSocket-Version: 13\r\n\r\n", okey)
	fmt.Println(send)
	con.Write([]byte(send))
	//con.
	handleClientRequest(con)
}
func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [2048]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
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
	fmt.Printf("read -> %v\nremote -> %s\n", string(b[:]), string(address))
	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Printf("链接请求")
		client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		//fmt.Fprint(client, )
	} else {
		fmt.Printf("链接问题")
		server.Write(b[:n])
	}

	defer server.Close()
	defer client.Close()
	go io.Copy(server, client)
	io.Copy(client, server)

}
