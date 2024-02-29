package router

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	controllers "nas/project/src/controllers"
	"nas/project/src/middleware"
	"strconv"
)

var NormalRouter = gin.Default()

func GetNormalRouter() *gin.Engine {
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	NormalRouter.Use(middleware.TlsHandler(serverPort))
	NormalRouter.GET("/hello", controllers.Test)
	NormalRouter.POST("/login", controllers.Login)
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
			userFileApi.GET("/large", controllers.LargeFileTransitionPrepare)
		}
		userDirApi := userApi.Group("/dir")
		{
			userDirApi.GET("/:catalog/:order/:page_num", controllers.CheckDir)
			userDirApi.POST("", controllers.CreateDir)
			userDirApi.DELETE("", controllers.DeleteDir)
		}
		userApi.GET("/thumbnail", controllers.GetThumbnail) //要改
	}
	/*adminRouter := NormalRouter.Group("/admin")*/
	return NormalRouter
}

func RunTLSOnConfig(router *gin.Engine) {
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	//router.Use(middleware.TlsHandler(serverPort))
	keyPath := Utils.DefaultConfigReader().Get("TLS:keyPath").(string)
	pemPath := Utils.DefaultConfigReader().Get("TLS:pemPath").(string)
	err := router.RunTLS(":"+strconv.Itoa(serverPort), pemPath, keyPath)
	if err != nil {
		return
	}
}
