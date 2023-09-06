package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"nas/project/controllers"
	"nas/project/src/Utils"
	"nas/project/src/middleware"
	"strconv"
)

var NormalRouter = gin.Default()

func GetNormalRouter() *gin.Engine {
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	NormalRouter.Use(cors.Default())
	NormalRouter.Use(middleware.TlsHandler(serverPort))
	NormalRouter.POST("/api/login", controllers.Login)
	//用户路由组
	userApi := NormalRouter.Group("/user")
	{
		userApi.Use(middleware.TokenInspect())
		userApi.GET("/info", controllers.GetUserInfo)
		userFileApi := userApi.Group("/file")
		{
			userFileApi.DELETE("/", controllers.DeleteFile)
			userFileApi.PUT("/", controllers.MoveFile)
			userFileApi.GET("/uploading", controllers.GetUnfinishedUpload)
			userFileApi.POST("/small", controllers.UploadSmallFile)
			userFileApi.POST("/large", controllers.UploadLargeFile)
		}
		userApi.GET("/dir/:dir_path/:order/:page_num", controllers.CheckDir)
	}
	/*adminRouter := NormalRouter.Group("/admin")*/
	return NormalRouter
}

func RunTLSOnConfig(router *gin.Engine) {
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	router.Use(middleware.TlsHandler(serverPort))
	keyPath := Utils.DefaultConfigReader().Get("TLS:keyPath").(string)
	pemPath := Utils.DefaultConfigReader().Get("TLS:pemPath").(string)
	err := router.RunTLS(":"+strconv.Itoa(serverPort), pemPath, keyPath)
	if err != nil {
		return
	}
}
