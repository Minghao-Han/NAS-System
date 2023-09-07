package main

import (
	"nas/project/src/Utils"
	"nas/project/src/router"
	"sync"
)

func main() {
	normalRouter := router.GetNormalRouter()
	csPorts := Utils.DefaultConfigReader().Get("FSP:csPorts").([]interface{})
	dsPorts := Utils.DefaultConfigReader().Get("FSP:dsPorts").([]interface{})
	wg := sync.WaitGroup{}
	wg.Add(len(csPorts) + len(dsPorts) + 1)
	for _, csPort := range csPorts {
		csPort := csPort
		csRouter := router.GetControlStreamRouter()
		go func() {
			router.RunOnConfig(csPort.(int), csRouter)
			wg.Done()
		}()
	}
	for _, dsPort := range dsPorts {
		dsPort := dsPort
		dsRouter := router.GetDataStreamRouter()
		go func() {
			router.RunOnConfig(dsPort.(int), dsRouter)
			wg.Done()
		}()
	}
	go func() {
		router.RunTLSOnConfig(normalRouter)
		wg.Done()
	}()
	wg.Wait()
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
