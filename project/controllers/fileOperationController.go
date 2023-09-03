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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "file doesn't exist",
			"err": existErr,
		})
		return
	}
	rmErr := os.Remove(filePath)
	if rmErr != nil { //删除文件失败
		c.JSON(http.StatusUnprocessableEntity, gin.H{
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
	SourceFilePath string `json:"sourceFilePath,omitempty"`
	DestinyPath    string `json:"destinyPath,omitempty"`
	Cover          bool   `json:"cover,omitempty"` //是否覆盖同名文件
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
	sourceFilePath := userRoot + mf.SourceFilePath
	destinyPath := userRoot + mf.DestinyPath
	fileInfo, existErr := os.Stat(sourceFilePath)
	if os.IsNotExist(existErr) { //文件不存在
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "file doesn't exist",
			"err": existErr,
		})
		return
	}
	if _, dirErr := os.Stat(destinyPath); os.IsNotExist(dirErr) { //目的文件夹不存在
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "destiny dir doesn't exist",
			"err": dirErr,
		})
		return
	}
	newFilePath := destinyPath + "/" + fileInfo.Name()
	_, err := os.Stat(newFilePath)
	if !os.IsNotExist(err) { //有同名文件
		if !mf.Cover { //不覆盖
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"msg": "Duplicate file",
				"err": err,
			})
			return
		}
	}
	renameErr := os.Rename(sourceFilePath, newFilePath)
	if renameErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "failed to move file",
			"err": renameErr,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "move file succeeded",
	})
	return
}
