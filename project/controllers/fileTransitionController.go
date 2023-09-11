package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/DA/userDA"
	"nas/project/src/Entities"
	"nas/project/src/Service"
	"nas/project/src/Service/PortManage"
	"nas/project/src/Utils"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	STOP        int64 = -1
	CANCEL      int64 = -2
	NO_OP       int64 = -3
	bufferSize        = int64(Utils.DefaultConfigReader().Get("download:bufferSize").(int))
	sectionNum        = Utils.DefaultConfigReader().Get("download:sectionNum").(int)
	sectionSize       = bufferSize / int64(sectionNum)
)
var portsManager = PortManage.DefaultPortsManager()

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

func LargeFileTransitionPrepare(c *gin.Context) { /*负责先做些基本的检查，包括margin，文件路径等。如果可以上传，则返回csPort,dsPort*/
	//path := c.GetHeader("path")
	//fileSize, err := strconv.ParseUint(c.GetHeader("size"), 10, 64)
	//filename := c.GetHeader("filename")
	//value, _ := c.Get("userId")
	//userId := value.(int)
	//if path == "" || filename == "" || err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"msg": "missing required info",
	//	})
	//	return
	//}
	//err = Service.LargeFileUploadPrepare(path, fileSize, userId)
	//if err != nil {
	//	c.JSON(http.StatusUnprocessableEntity, gin.H{
	//		"msg": err.Error(),
	//	})
	//	return
	//}
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

func CsForUpload(c *gin.Context) {
	///*验证connection预留信息*/
	//sourceIP := net.ParseIP(c.ClientIP())
	//_, portStr, _ := net.SplitHostPort(c.Request.Host)
	//port, err := strconv.Atoi(portStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"msg": err.Error(),
	//	})
	//	return
	//}
	//portEntity, found := portsManager.FindPort(0, port)
	//if !found {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"msg": "didn't find port info",
	//	})
	//	return
	//}
	//connection, reserved := portEntity.FindConnection(sourceIP)
	//if !reserved {
	//	c.JSON(http.StatusUnprocessableEntity, gin.H{
	//		"msg": "reservation didn't find",
	//	})
	//	return
	//}
	//defer portEntity.DisConnectByIP(sourceIP) //断开连接
	//connection.GetDS2CS() <- -1
	//cs2ds := connection.GetCS2DS()

}

// DsForUpload header: token,offset,filename,path,fileSize,uploadId
func DsForUpload(c *gin.Context) {
	/*获取和检查参数*/
	/*验证margin*/
	/*开始写*/
	/*检查连接*/
	port, connection, err := verifyConnReservation(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer port.DisConnectByIP(net.ParseIP(c.ClientIP()))
	connection.GetDS2CS() <- -1
	/*获取参数*/
	offset, err := strconv.Atoi(c.GetHeader("offset"))
	filename := c.GetHeader("filename")
	clientFilePath := c.GetHeader("clientFilePath")
	path := c.GetHeader("path")
	filePath := path + "/" + filename
	fileSize, err := strconv.Atoi(c.GetHeader("fileSize"))
	value, _ := c.Get("userId")
	uploadId, err := strconv.Atoi(c.GetHeader("uploadId"))
	userId := value.(int)
	user, err := userDA.FindById(userId)
	if path == "" || filename == "" || clientFilePath == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing required info \n" + err.Error(),
		})
		return
	}
	if uint64(fileSize-offset) > user.Margin {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "no more space for this file",
		})
		return
	}
	fullFilePath := Service.GetFullFilePath(filePath, userId)
	/*如果是续传就检查数据库中的项正不正确，如果是新上传就要检查是否有同名文件*/
	var file *os.File
	if uploadId != -1 {
		uploadLog, _ := uploadDA.FindById(uploadId)
		if uploadLog.Finished == true || uploadLog.Received_bytes != uint64(offset) || uploadLog.Path != filePath {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"msg": "resume upload failed",
			})
		}
		file, err = os.Open(fullFilePath)
		defer file.Close()
	} else {
		fullFilePath = Service.DuplicateFileName(fullFilePath)
		uploadId, _ = uploadDA.Insert(Entities.UploadLog{
			Id:             0,
			Uploader:       userId,
			Path:           Service.GetUserRelativePath(fullFilePath, userId),
			Finished:       false,
			Received_bytes: 0,
			Size:           uint64(fileSize),
			ClientFilePath: clientFilePath,
		})
		file, err = os.Create(fullFilePath)
		defer file.Close()
	}
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	/*从请求头读取*/
	plaintext := make([]byte, bufferSize)
	ciphertext := make([]byte, bufferSize)
	chaDecipher := Utils.DefaultChaEncryptor()
	var receivedBytes uint64 = 0
	requestBodyReader := c.Request.Body
	for {
		n, readErr := requestBodyReader.Read(plaintext)
		if readErr != nil && readErr != io.EOF && !errors.Is(readErr, http.ErrBodyReadAfterClose) { //读取出错
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"msg": readErr.Error(),
			})
			break
		}
		//if n == 0 && (readErr == io.EOF || errors.Is(readErr, http.ErrBodyReadAfterClose)) { //读到结束或连接断开
		if n == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "uploaded",
			})
			break
		}
		loopTimes := int((int64(n) + sectionSize - 1) / sectionSize)
		for i := 0; i < loopTimes; i++ {
			index := int64(i)
			upto := min((index+1)*sectionSize, int64(n))
			wg := sync.WaitGroup{}
			wg.Add(loopTimes)
			go func() {
				chaDecipher.Encrypt(plaintext[index*sectionSize:upto], ciphertext[index*sectionSize:upto])
				wg.Done()
			}()
			wg.Wait()
		}
		if _, err = file.Write(ciphertext[:n]); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"msg": err.Error(),
			})
			break
		}
		receivedBytes += uint64(n)
	}
	uploadLog, _ := uploadDA.FindById(uploadId)
	uploadLog.Received_bytes = receivedBytes
	uploadDA.Update(*uploadLog)
	return
}

func CsForDownload(c *gin.Context) {
	var opcode int64
	/**/
	err := c.ShouldBindJSON(&opcode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	_, connection, err := verifyConnReservation(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	cs2ds := connection.GetCS2DS()
	go func() {
		cs2ds <- opcode
	}()
	return
}

func DsForDownload(c *gin.Context) {
	/*验证connection预留信息*/
	port, connection, err := verifyConnReservation(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer port.DisConnectByIP(net.ParseIP(c.ClientIP())) //断开连接
	cs2ds := connection.GetCS2DS()
	/*获取请求头其他信息*/
	value, _ := c.Get("userId")
	userId := value.(int)
	filePath := c.GetHeader("filePath")
	/*获取文件reader*/
	file, err := os.Open(Service.GetFullFilePath(filePath, userId))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	connection.GetDS2CS() <- -1
	defer file.Close()
	///*设置响应头*/
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Connection", "close")
	///*写入输出流*/
	chaDecipher := Utils.DefaultChaEncryptor()
	var offset int64 = 0
	ciphertext := make([]byte, bufferSize)
	plaintext := make([]byte, bufferSize)
	c.Stream(func(w io.Writer) bool {
		for {
			/*检查opcode*/
			var opcode int64 = NO_OP
			select {
			case opcode = <-cs2ds:
			default:
				break
			}
			if opcode == STOP {
				break
			} else if opcode >= 0 { //opcode为正数时表示偏移多少
				offset = opcode
			} else if opcode == NO_OP {
			} else { //错误的opcode
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "opcode invalid",
				})
				break
			}
			if c.IsAborted() {
				c.AbortWithStatus(499)
				break
			}
			/*解密*/
			n, err := file.ReadAt(ciphertext, offset)
			if err != nil && err.Error() != "EOF" { //文件读取出错
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": err.Error(),
				})
				break
			}
			if n == 0 {
				break //跳出for循环
			}
			//向上取整
			loopTimes := int((int64(n) + sectionSize - 1) / sectionSize)
			wg := sync.WaitGroup{}
			wg.Add(int(loopTimes))
			for index := 0; index < loopTimes; index++ {
				index := int64(index)
				upto := min((index+1)*sectionSize, int64(n))
				go func() {
					chaDecipher.Decrypt(ciphertext[index*sectionSize:upto], plaintext[index*sectionSize:upto])
					wg.Done()
				}()
			}
			wg.Wait()
			//w.Write(plaintext[:n])
			w.Write(ciphertext[:n])
			offset += int64(n)
		}
		return false
	})
	//c.Writer.Flush()
	return
}

//	func DsForDownloadTest(c *gin.Context) {
//		c.Header("Content-Disposition", "attachment; filename=aa")
//		c.Header("Content-Type", "application/octet-stream")
//		c.Header("Connection", "close")
//		for i := 0; i < 3; i++ {
//			c.Stream(func(w io.Writer) bool {
//				w.Write([]byte(time.Now().String()))
//				w.Write([]byte("hello world" + strconv.Itoa(i) + "\n"))
//				return false
//			})
//		}
//	}
func verifyConnReservation(c *gin.Context) (*PortManage.Port, *PortManage.Connection, error) {
	sourceIP := net.ParseIP(c.ClientIP())
	_, portStr, _ := net.SplitHostPort(c.Request.Host)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, nil, err
	}
	portEntity, found := portsManager.FindPort(0, port)
	if !found {
		return nil, nil, fmt.Errorf("didn't find port info")
	}
	connection, reserved := portEntity.FindConnection(sourceIP)
	if !reserved {
		return nil, nil, fmt.Errorf("reservation didn't find")
	}
	return portEntity, connection, nil
}
