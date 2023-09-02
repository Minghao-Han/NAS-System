package main

import (
	"nas/project/src/router"
)

func main() {
	//pw := []byte("123456")
	//cryptespw, _ := Utils.DefaultAESEncryptor().EncryptWithPadding(pw)
	//hmh := Entities.User{
	//	UserId:   0,
	//	UserName: "hmh",
	//	Password: cryptespw,
	//	Capacity: 150,
	//	Margin:   88,
	//}
	//userDA.Insert(hmh)
	//GinHttps(true) // 这里false 表示 http 服务，非 https
	normalRouter := router.GetNormalRouter()
	router.RunTLSOnConfig(normalRouter)
}

//
//func GinHttps(isHttps bool) error {
//
//	r := gin.Default()
//	r.GET("/test", func(c *gin.Context) {
//		c.String(200, "test for 【%s】", "https")
//	})
//	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
//	if isHttps {
//		r.Use(middleware.TlsHandler(serverPort))
//		keyPath := Utils.DefaultConfigReader().Get("TLS:keyPath").(string)
//		pemPath := Utils.DefaultConfigReader().Get("TLS:pemPath").(string)
//		return r.RunTLS(":"+strconv.Itoa(serverPort), pemPath, keyPath)
//	}
//
//	return r.Run(":" + strconv.Itoa(serverPort))
//}
