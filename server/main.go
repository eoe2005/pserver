package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/eoe2005/pserver"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{}

func main() {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8888"
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": "",
		})
	})
	r.GET("/http", func(c *gin.Context) {
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()
	})

	r.GET("/sock5", func(c *gin.Context) {
		pserver.Debug("进入sock5模式")
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			pserver.Debug("建立ws失败 %s", err.Error())
			return
		}
		pserver.Sock5Raw(ws)
		pserver.Debug("sock5结束")
	})
	r.GET("/wssock5", func(c *gin.Context) {
		pserver.Debug("进入sock5模式")

		h, ok := c.Writer.(http.Hijacker)
		if ok {
			con, _, e := h.Hijack()
			if e == nil {
				con.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: R0pRjt+Rr25+h9Irjy+Kx1ltJ0A=\r\n"))
				helo(con)
			} else {
				pserver.Debug("获取链接失败 %s", e.Error())
			}
		} else {
			pserver.Debug("转换失败")
		}

	})
	pserver.Debug("服务开始运行")
	r.Run(":" + p)
}
func helo(con net.Conn) {
	h1 := make([]byte, 3)
	con.Read(h1)
	con.Write([]byte{0x05, 0x00})
	h2 := make([]byte, 4)
	con.Read(h2)

	addr := ""
	con.Write([]byte{0x05, 0x00, 0x00, h2[3]})
	switch h2[3] {
	case 0x01:
		v := make([]byte, 4)
		con.Read(v)
		addr = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
		fmt.Printf("链接地址 %s\n", addr)
		con.Write(v)
	case 0x04:
		v := make([]byte, 16)
		con.Read(v)
		con.Write(v)
	case 0x03:
		l := make([]byte, 1)
		con.Read(l)
		v := make([]byte, l[0])
		con.Read(v)
		con.Write(l)
		con.Write(v)
		addr = string(v)
		fmt.Printf("链接地址 %s\n", addr)
	}
	if addr == "" {
		con.Close()
	}
	p := make([]byte, 2)
	con.Read(p)
	port := (int(p[0]) << 8) + int(p[1])
	fmt.Printf("链接端口 %d - %d %d %s\n", port, int(p[0])<<8, p[1], string(p))
	desc, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	con.Write(p)
	go trance(con, desc)
}

// / 数据交换
func trance(src, desc net.Conn) {
	defer src.Close()
	defer desc.Close()
	go io.Copy(src, desc)
	io.Copy(desc, src)

}
