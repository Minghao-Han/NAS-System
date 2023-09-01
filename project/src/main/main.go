package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

//func main() {
//	//portManage := PortManage.NewPortsManager()
//	//for i := 0; i < 2; i++ {
//	//	go func() {
//	//		csPort, dsPort, connIndex, ok := portManage.PrepareConnection(net.ParseIP("192.168.1.1"))
//	//		fmt.Printf("csPort is %d,dsPort is %d,connIndex is %d,ok is %t \n", csPort, dsPort, connIndex, ok)
//	//	}()
//	//}
//	//time.Sleep(5 * time.Second)
//	//return
//}
import (
	"github.com/unrolled/secure"
)

func main() {
	GinHttps(false) // 这里false 表示 http 服务，非 https
}

func GinHttps(isHttps bool) error {

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test for 【%s】", "https")
	})

	if isHttps {
		r.Use(TlsHandler(8000))

		return r.RunTLS(":"+strconv.Itoa(8000), "/path/to/test.pem", "/path/to/test.key")
	}

	return r.Run(":" + strconv.Itoa(8000))
}

func TlsHandler(port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
