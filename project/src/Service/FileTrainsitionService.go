package Service

import (
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
	RESUME      int = -1
	CANCEL      int = -2
	NoOp        int = -3
	bufferSize      = int64(Utils.DefaultConfigReader().Get("download:bufferSize").(int))
	sectionNum      = Utils.DefaultConfigReader().Get("download:sectionNum").(int)
	sectionSize     = bufferSize / int64(sectionNum)
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
		Id:            0,
		Uploader:      userId,
		Path:          strings.ReplaceAll(filePath, userRoot, ""),
		Finished:      true,
		ReceivedBytes: fileSize,
		Size:          fileSize,
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
	var receivedBytes uint64 = 0
	if uploadId == RESUME {
		uploadLog, _ := uploadDA.FindById(uploadId)
		if uploadLog.Finished == true || uploadLog.ReceivedBytes != offset || uploadLog.Path != uploadPath {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"msg": "resume upload failed",
			})
		}
		receivedBytes = uploadLog.ReceivedBytes
		file, err = os.Open(fullFilePath)
		defer file.Close()
	} else {
		fullFilePath = DuplicateFileName(fullFilePath)
		uploadId, _ = uploadDA.Insert(Entities.UploadLog{
			Id:             0,
			Uploader:       userId,
			Path:           GetUserRelativePath(fullFilePath, userId),
			Finished:       false,
			ReceivedBytes:  24, //预留24，因为每个文件都要使用随机生成的nonce加密，nonce为24B，并明文放在文件头
			Size:           fileSize,
			ClientFilePath: clientFilePath,
		})
		file, err = os.Create(fullFilePath)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}(file)
	}
	if err != nil {
		return err
	}
	cha20FileIO, err := Utils.DefaultChaCha20FileIO(file, file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	plaintext := make([]byte, bufferSize)
	/*从请求头读取*/
	/*read from request header*/
	fileSize -= receivedBytes //减掉已上传的部分
	receivedBytes = 0         //receivedBytes重新置0，表示这次接收的字节数
	for {
		// http 长连接，可以不断从body中读。
		n, readErr := c.Request.Body.Read(plaintext)
		if readErr != nil && readErr != io.EOF && !errors.Is(readErr, http.ErrBodyReadAfterClose) { //读取出错 read error
			break
		}
		if n == 0 || receivedBytes >= fileSize {
			break
		}

		//使用chacha20 代理的writer
		if _, err := cha20FileIO.Write(plaintext[:n], file); err != nil { //写入错误
			break
		}
		receivedBytes += uint64(n)
	}
	err = cha20FileIO.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "uploaded",
	})
	//更新上传记录
	//update upload log
	uploadLog, _ := uploadDA.FindById(uploadId)
	uploadLog.ReceivedBytes += receivedBytes
	_, err = uploadDA.Update(*uploadLog)
	if err != nil {
		return err
	}
	//更新用户容量
	//update user's margin
	user, err := userDA.FindById(userId)
	user.Margin -= receivedBytes
	_, err = userDA.Update(*user)
	if err != nil {
		return err
	}
	return nil
}
