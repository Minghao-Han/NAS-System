package router

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/middleware"
)

func DefaultCorsRouter() *gin.Engine {
	defaultRouter := gin.Default()
	defaultRouter.Use(middleware.CORSMiddleware())
	return defaultRouter
}
