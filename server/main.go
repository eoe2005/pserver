package main

import (
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
		pserver.Sock5(ws)
		pserver.Debug("sock5结束")
	})
	pserver.Debug("服务开始运行")
	r.Run(":" + p)
}
