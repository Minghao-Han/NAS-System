package controllers

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Service"
	"net/http"
	"strconv"
	"strings"
)

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
	err := Service.DeleteFile(userId, df.FilePath)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "delete file succeeded",
		})
	}
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
	err := Service.MoveFile(userId, mf.SourceFilePath, mf.DestinyPath, mf.Cover)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "move file succeeded",
		})
	}
	return
}

// CheckDir url dir/:dir_path/:style/:order/:page_num
func CheckDir(c *gin.Context) {
	value, _ := c.Get("userId")
	userId := value.(int)
	dirPath := c.Param("dir_path")
	dirPath = strings.ReplaceAll(dirPath, "_", "/")
	order := c.Param("order")
	pageNum, err := strconv.Atoi(c.Param("page_num"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	/*
		从redis中取，后面迭代
	*/
	if dirPath == "" || order == "" { //缺少参数
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "require params",
		})
		return
	}
	fileInfos, checkDirErr := Service.CheckDir(userId, dirPath, order, pageNum)
	if checkDirErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": checkDirErr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":       "query ok",
		"fileInfos": fileInfos,
	})
}
