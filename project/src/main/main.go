package main

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	"nas/project/src/middleware"
	"strconv"
)

func main() {
	//GinHttps(true) // 这里false 表示 http 服务，非 https
}

func GinHttps(isHttps bool) error {

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test for 【%s】", "https")
	})
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	if isHttps {
		r.Use(middleware.TlsHandler(serverPort))
		keyPath := Utils.DefaultConfigReader().Get("TLS:keyPath").(string)
		pemPath := Utils.DefaultConfigReader().Get("TLS:pemPath").(string)
		return r.RunTLS(":"+strconv.Itoa(serverPort), pemPath, keyPath)
	}

	return r.Run(":" + strconv.Itoa(serverPort))
}
