package router

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	"nas/project/src/controllers"
	"nas/project/src/middleware"
	"strconv"
)

func GetControlStreamRouter() *gin.Engine {
	var CSRouter = gin.Default()
	CSRouter.Use(middleware.TokenInspect())
	CSRouter.POST("/download", controllers.CsForDownload)
	CSRouter.POST("/upload", controllers.CsForUpload)
	return CSRouter
}
func GetDataStreamRouter() *gin.Engine {
	var DSRouter = gin.Default()
	DSRouter.Use(middleware.TokenInspect())
	DSRouter.GET("/download", controllers.DsForDownload)
	DSRouter.POST("/upload", controllers.DsForUpload)
	return DSRouter
}

func RunOnConfig(port int, router *gin.Engine) {
	router.Use(middleware.TlsHandler(port))
	keyPath := Utils.DefaultConfigReader().Get("TLS:keyPath").(string)
	pemPath := Utils.DefaultConfigReader().Get("TLS:pemPath").(string)
	err := router.RunTLS(":"+strconv.Itoa(port), pemPath, keyPath)
	if err != nil {
		return
	}
}
