package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
	l, e := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if e != nil {
		panic("监听端口失败")
	}
	for {
		con, e := l.Accept()
		if e != nil {
			fmt.Printf("链接数据失败 %s\n", e.Error())
		} else {
			go handleClientRequest(con)
		}
	}
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
