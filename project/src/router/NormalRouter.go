package router

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	controllers2 "nas/project/src/controllers"
	"nas/project/src/middleware"
	"net/http"
	"strconv"
)

var NormalRouter = gin.Default()

func GetNormalRouter() *gin.Engine {
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	NormalRouter.Use(middleware.TlsHandler(serverPort))
	NormalRouter.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello from hmh")
		return
	})
	NormalRouter.POST("/login", controllers2.Login)
	//用户路由组
	userApi := NormalRouter.Group("/user")
	{
		userApi.Use(middleware.TokenInspect())
		userApi.GET("/info", controllers2.GetUserInfo)
		userFileApi := userApi.Group("/file")
		{
			userFileApi.DELETE("/", controllers2.DeleteFile)
			userFileApi.PUT("/", controllers2.MoveFile)
			userFileApi.GET("/uploading", controllers2.GetUnfinishedUpload)
			userFileApi.POST("/small", controllers2.UploadSmallFile)
			userFileApi.GET("/large", controllers2.LargeFileTransitionPrepare)
		}
		userDirApi := userApi.Group("/dir")
		{
			userDirApi.GET("/:dir_path/:catalog/:order/:page_num", controllers2.CheckDir)
			userDirApi.POST("", controllers2.CreateDir)
			userDirApi.DELETE("", controllers2.DeleteDir)
		}
		userApi.GET("/thumbnail/:file_path", controllers2.GetThumbnail)
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
