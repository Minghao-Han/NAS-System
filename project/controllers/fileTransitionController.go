package controllers

import (
	"github.com/gin-gonic/gin"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/Service"
	"nas/project/src/Service/PortManage"
	"net"
	"net/http"
	"strconv"
)

func GetUnfinishedUpload(c *gin.Context) {
	value, _ := c.Get("userId")
	userId := value.(int)
	unfinishedUpload, err := uploadDA.FindUnfinishedByUploader(strconv.Itoa(userId))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":               "query ok",
		"unfinished_upload": unfinishedUpload,
	})
	return
}

func UploadSmallFile(c *gin.Context) {
	path := c.GetHeader("path")
	fileSize, err := strconv.ParseUint(c.GetHeader("size"), 10, 64)
	filename := c.GetHeader("filename")
	value, _ := c.Get("userId")
	userId := value.(int)
	if path == "" || filename == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing required info",
		})
		return
	}
	if fileSize > (50 * 1024 * 1024) { //100m
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"msg": "file too large, use get_port for large file upload",
		})
		return
	}
	err = Service.UploadSmallFile(c, path, fileSize, filename, userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "upload file succeeded",
	})
}

func UploadLargeFile(c *gin.Context) { /*负责先做些基本的检查，包括margin，文件路径等。如果可以上传，则返回csPort,dsPort*/
	path := c.GetHeader("path")
	fileSize, err := strconv.ParseUint(c.GetHeader("size"), 10, 64)
	filename := c.GetHeader("filename")
	value, _ := c.Get("userId")
	userId := value.(int)
	if path == "" || filename == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing required info",
		})
		return
	}
	err = Service.LargeFileUploadPrepare(path, fileSize, userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//初步判断可以上传，返回csPort,dsPort
	csPort, dsPort, connIndex, got := PortManage.DefaultPortsManager().PrepareConnection(net.ParseIP(c.ClientIP()))
	if !got { //没能预留连接
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "failed to get connection reservation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "connection reserve succeeded",
		"csPort":    csPort,
		"dsPort":    dsPort,
		"connIndex": connIndex,
	})
	return
}
