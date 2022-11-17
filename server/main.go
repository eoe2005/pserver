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
	_, e := con.Read(h1)
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}

	_, e = con.Write([]byte{0x05, 0x00})
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}

	h2 := make([]byte, 4)
	_, e = con.Read(h2)
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}

	addr := ""
	_, e = con.Write([]byte{0x05, 0x00, 0x00, h2[3]})
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}

	switch h2[3] {
	case 0x01:
		v := make([]byte, 4)
		_, e = con.Read(v)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}

		addr = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
		fmt.Printf("链接地址 %s\n", addr)
		_, e = con.Write(v)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}

	case 0x04:
		v := make([]byte, 16)
		_, e = con.Read(v)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}

		_, e = con.Write(v)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}
	case 0x03:
		l := make([]byte, 1)
		_, e = con.Read(l)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}
		v := make([]byte, l[0])
		_, e = con.Read(v)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}

		_, e = con.Write(l)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}

		_, e = con.Write(v)
		if e != nil {
			pserver.Debug("sock init err %s", e.Error())
			return
		}
		addr = string(v)
		fmt.Printf("链接地址 %s\n", addr)
	}
	if addr == "" {
		con.Close()
	}
	p := make([]byte, 2)
	_, e = con.Read(p)
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}

	port := (int(p[0]) << 8) + int(p[1])
	fmt.Printf("链接端口 %d - %d %d %s\n", port, int(p[0])<<8, p[1], string(p))
	desc, e := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}

	_, e = con.Write(p)
	if e != nil {
		pserver.Debug("sock init err %s", e.Error())
		return
	}
	go trance(con, desc)
}

// / 数据交换
func trance(src, desc net.Conn) {
	defer src.Close()
	defer desc.Close()
	go io.Copy(src, desc)
	io.Copy(desc, src)

}
