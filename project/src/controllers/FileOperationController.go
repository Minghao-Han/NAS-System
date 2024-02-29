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
	userId := GetUserIdFromContext(c)
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
	userId := GetUserIdFromContext(c)
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

// CheckDir url /:dir_path/:catalog/:order/:page_num
func CheckDir(c *gin.Context) {
	value, _ := c.Get("userId")
	userId := value.(int)
	dirPath := c.GetHeader("dir_path")
	dirPath = strings.ReplaceAll(dirPath, "_", Service.Slash)
	order := c.Param("order")
	catalog := c.Param("catalog")
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
	fileInfos, checkDirErr := Service.CheckDir(userId, dirPath, order, pageNum, catalog)
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

type NewDirInfo struct {
	PrentDirPath string `json:"parent_dir_path,omitempty"`
	DirName      string `json:"dir_name,omitempty"`
}

func CreateDir(c *gin.Context) {
	userId := GetUserIdFromContext(c)
	ndi := NewDirInfo{}
	if err := c.ShouldBindJSON(&ndi); err != nil || len(ndi.PrentDirPath) == 0 || len(ndi.DirName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing required params",
		})
		return
	}
	err := Service.CreateDir(ndi.PrentDirPath, ndi.DirName, userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "create folder succeed",
	})
}

type deleteInfo struct {
	DirPath string `json:"dir_path,omitempty"`
}

func DeleteDir(c *gin.Context) {
	userId := GetUserIdFromContext(c)
	deleteInfo := deleteInfo{}
	if err := c.ShouldBindJSON(&deleteInfo); err != nil || len(deleteInfo.DirPath) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing required params",
		})
		return
	}
	if err := Service.DeleteDir(deleteInfo.DirPath, userId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "delete  folder succeed",
	})
}

// GetThumbnail file_path
func GetThumbnail(c *gin.Context) {
	filePath := c.GetHeader("path")
	userId := GetUserIdFromContext(c)
	thumbnailBuf, err := Service.GetThumbnail(userId, filePath)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//c.Header("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)
	thumbnailBuf.WriteTo(c.Writer)
	return
}

func GetUserIdFromContext(c *gin.Context) int {
	value, _ := c.Get("userId")
	return value.(int)
}
