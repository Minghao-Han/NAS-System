package router

import (
	"github.com/gin-gonic/gin"
	"nas/project/controllers"
	"nas/project/src/Utils"
	"nas/project/src/middleware"
	"strconv"
)

var NormalRouter = gin.Default()

func GetNormalRouter() *gin.Engine {
	serverPort := Utils.DefaultConfigReader().Get("Server:port").(int)
	NormalRouter.Use(middleware.TlsHandler(serverPort))
	NormalRouter.POST("/login", controllers.Login)
	//用户路由组
	userRouter := NormalRouter.Group("/user")
	{
		userRouter.Use(middleware.TokenInspect())
		userRouter.GET("/info", controllers.GetUserInfo)
		userRouter.DELETE("/file", controllers.DeleteFile)
	}
	//adminRouter := NormalRouter.Group("/admin")
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
