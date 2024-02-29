package Service

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/DA/userDA"
	"nas/project/src/Entities"
	"nas/project/src/Utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	STOP        int64 = -1
	CANCEL      int64 = -2
	NoOp        int64 = -3
	bufferSize        = int64(Utils.DefaultConfigReader().Get("download:bufferSize").(int))
	sectionNum        = Utils.DefaultConfigReader().Get("download:sectionNum").(int)
	sectionSize       = bufferSize / int64(sectionNum)
)

// UploadSmallFile call this function to upload a small file.
func UploadSmallFile(c *gin.Context, path string, fileSize uint64, fileName string, userId int) error {
	err := MarginAvailable(userId, fileSize)
	if err != nil {
		return err
	}
	// 获取文件字节流
	fileData, err := c.GetRawData()
	if err != nil {
		return err
	}
	/*
		bytes 用chacha20加密
	*/
	cipherData := make([]byte, len(fileData))
	encryptErr := Utils.DefaultChaEncryptor().Encrypt(fileData, cipherData)
	if encryptErr != nil {
		return encryptErr
	}
	destinyPath := diskRoot + strconv.Itoa(userId) + path
	_, dirErr := os.Stat(destinyPath)

	if os.IsNotExist(dirErr) {
		return fmt.Errorf("destiny dir doesn't exist")
	}
	/*解决文件重名问题*/
	/*Solve the issue of duplicate file names.*/
	filePath := diskRoot + strconv.Itoa(userId) + path + "/" + fileName
	_, duplicateErr := os.Stat(filePath)
	index := 0
	newFilePath := filePath
	for !os.IsNotExist(duplicateErr) {
		newFilePath = filePath
		index++
		ext := filepath.Ext(newFilePath)
		base := strings.TrimSuffix(newFilePath, ext)

		newBase := base + "_" + strconv.Itoa(index)
		newFilePath = newBase + ext
		_, duplicateErr = os.Stat(newFilePath)
	}
	filePath = newFilePath
	writeErr := os.WriteFile(filePath, cipherData, 0666)
	if err != nil {
		return writeErr
	}
	userRoot := diskRoot + strconv.Itoa(userId)
	_, daErr := uploadDA.Insert(Entities.UploadLog{
		Id:             0,
		Uploader:       userId,
		Path:           strings.ReplaceAll(filePath, userRoot, ""),
		Finished:       true,
		Received_bytes: fileSize,
		Size:           fileSize,
	})
	if daErr != nil {
		fmt.Println(daErr.Error())
	}
	return nil
}

// LargeFileUploadPrepare preparation is needed for large file upload. This function will inspect whether there is enough space for the file and verify the file path
func LargeFileUploadPrepare(path string, fileSize uint64, userId int) error {
	err := MarginAvailable(userId, fileSize)
	if err != nil {
		return err
	}
	//目的目录检查
	destinyPath := diskRoot + strconv.Itoa(userId) + path
	_, dirErr := os.Stat(destinyPath)
	if os.IsNotExist(dirErr) {
		fmt.Println(destinyPath)
		return fmt.Errorf("destiny folder doesn't exist")
	}
	return nil
}

// DuplicateFileName check whether there is a file of which the path is the same as the that of new file's
// If there's a namesake, add _1 _2 ... after the file name as suffix and then return.
func DuplicateFileName(filePath string) string {
	_, duplicateErr := os.Stat(filePath)
	index := 0
	newFilePath := filePath
	for !os.IsNotExist(duplicateErr) {
		newFilePath = filePath
		index++
		ext := filepath.Ext(newFilePath)
		base := strings.TrimSuffix(newFilePath, ext)

		newBase := base + "_" + strconv.Itoa(index)
		newFilePath = newBase + ext
		_, duplicateErr = os.Stat(newFilePath)
	}
	filePath = newFilePath
	return filePath
}

// Upload upload large file
func Upload(c *gin.Context, offset uint64, uploadPath string, fileSize uint64, uploadId int, userId int, clientFilePath string) error {
	err := MarginAvailable(userId, fileSize-offset)
	if err != nil {
		return fmt.Errorf("no more space for this file")
	}
	fullFilePath := GetFullFilePath(uploadPath, userId)
	/*如果是续传就检查数据库中的项正不正确，如果是新上传就要检查是否有同名文件*/
	/*If it's a resume upload, check if the items in the database are correct. If it's a new upload, check for the existence of a file with the same name.*/
	var file *os.File
	if uploadId != -1 {
		uploadLog, _ := uploadDA.FindById(uploadId)
		if uploadLog.Finished == true || uploadLog.Received_bytes != offset || uploadLog.Path != uploadPath {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"msg": "resume upload failed",
			})
		}
		file, err = os.Open(fullFilePath)
		defer file.Close()
	} else {
		fullFilePath = DuplicateFileName(fullFilePath)
		uploadId, _ = uploadDA.Insert(Entities.UploadLog{
			Id:             0,
			Uploader:       userId,
			Path:           GetUserRelativePath(fullFilePath, userId),
			Finished:       false,
			Received_bytes: 0,
			Size:           fileSize,
			ClientFilePath: clientFilePath,
		})
		file, err = os.Create(fullFilePath)
		defer file.Close()
	}
	if err != nil {
		return err
	}
	var receivedBytes uint64 = 0
	fileWriter := bufio.NewWriter(file)
	plaintext := make([]byte, bufferSize)
	/*从请求头读取*/
	/*read from request header*/
	for {
		n, readErr := c.Request.Body.Read(plaintext)
		if receivedBytes >= 9752518000 {
			fmt.Println("here")
		}
		if readErr != nil && readErr != io.EOF && !errors.Is(readErr, http.ErrBodyReadAfterClose) { //读取出错
			break
		}
		if n == 0 || receivedBytes >= fileSize {
			break
		}
		if _, err := fileWriter.Write(plaintext[:n]); err != nil { //写入错误
			break
		}
		receivedBytes += uint64(n)
	}
	fileWriter.Flush()
	c.JSON(http.StatusOK, gin.H{
		"msg": "uploaded",
	})
	//更新上传记录
	//update upload log
	uploadLog, _ := uploadDA.FindById(uploadId)
	uploadLog.Received_bytes = receivedBytes
	uploadDA.Update(*uploadLog)
	//更新用户容量
	//update user's margin
	user, err := userDA.FindById(userId)
	user.Margin -= receivedBytes
	userDA.Update(*user)
	return nil
}
