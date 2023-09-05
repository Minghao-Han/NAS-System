package Service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/Entities"
	"nas/project/src/Utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

func LargeFileUploadPrepare(path string, fileSize uint64, userId int) error {
	err := MarginAvailable(userId, fileSize)
	if err != nil {
		return err
	}
	//目的目录检查
	destinyPath := diskRoot + strconv.Itoa(userId) + path
	_, dirErr := os.Stat(destinyPath)
	if os.IsNotExist(dirErr) {
		return fmt.Errorf("destiny dir doesn't exist")
	}
	return nil
}
