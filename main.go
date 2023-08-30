package main

import (
	"github.com/gin-gonic/gin"
	"nas/app/utils"
	"net/http"
)

func main() {
	//创建一个路由Handler
	router := gin.Default()

	//get方法的查询
	router.GET("/query", func(c *gin.Context) {
		utils.Query()
		c.JSON(http.StatusOK, "AAAAAA")
	})
	router.Run(":8000")
}
