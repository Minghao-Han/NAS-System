package controllers

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	"net/http"
	"os"
	"strconv"
)

var diskRoot = Utils.DefaultConfigReader().Get("DiskRoot").(string)

/*
删除文件
*/
type deleteFile struct {
	FilePath string `json:"filePath,omitempty"`
}

func DeleteFile(c *gin.Context) {
	value, _ := c.Get("userId")
	userId := value.(int)
	var df = deleteFile{}
	if err := c.ShouldBindJSON(&df); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	filePath := diskRoot + strconv.Itoa(userId) + df.FilePath
	_, existErr := os.Stat(filePath)
	if os.IsNotExist(existErr) { //文件不存在
		c.JSON(http.StatusNoContent, gin.H{
			"msg": "file doesn't exist",
			"err": existErr,
		})
		return
	}
	rmErr := os.Remove(filePath)
	if rmErr != nil { //删除文件失败
		c.JSON(http.StatusNoContent, gin.H{
			"msg": "failed to delete file",
			"err": rmErr,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully removed file",
	})
	return
}

/*
移动文件
*/
type moveFile struct {
	FilePath    string `json:"filePath,omitempty"`
	DestinyPath string `json:"destinyPath,omitempty"`
	Cover       bool   `json:"cover,omitempty"` //是否覆盖同名文件
}

func MoveFile(c *gin.Context) {
	value, _ := c.Get("userId")
	userId := value.(int)
	var mf = moveFile{}
	if err := c.ShouldBindJSON(&mf); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	userRoot := diskRoot + strconv.Itoa(userId)
	filePath := userRoot + mf.FilePath
	destinyPath := userRoot + mf.DestinyPath
	_, existErr := os.Stat(filePath)
	if os.IsNotExist(existErr) { //文件不存在
		c.JSON(http.StatusNoContent, gin.H{
			"msg": "file doesn't exist",
			"err": existErr,
		})
		return
	}
	if _, dirErr := os.Stat(destinyPath); os.IsNotExist(dirErr) { //目的文件夹不存在
		c.JSON(http.StatusNoContent, gin.H{
			"msg": "destiny dir doesn't exist",
			"err": dirErr,
		})
		return
	}
	newFilePath := destinyPath+"/"+
}
