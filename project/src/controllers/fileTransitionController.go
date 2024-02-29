package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/Service"
	"nas/project/src/Service/PortManage"
	"nas/project/src/Utils"
	"net"
	"net/http"
	"os"
	"strconv"
)

var (
	STOP        int64 = -1
	CANCEL      int64 = -2
	NO_OP       int64 = -3
	bufferSize        = int64(Utils.DefaultConfigReader().Get("Download:bufferSize").(int))
	sectionNum        = Utils.DefaultConfigReader().Get("Download:sectionNum").(int)
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

// LargeFileTransitionPrepare 负责先做些基本的检查，包括margin，文件路径等。如果可以上传，则返回csPort,dsPort
func LargeFileTransitionPrepare(c *gin.Context) {
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
	prepareErr := Service.LargeFileUploadPrepare(path, fileSize, userId)
	if prepareErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": prepareErr.Error(),
		})
	}
	//初步判断可以上传，返回csPort,dsPort
	csPort, dsPort, connIndex, got := PortManage.DefaultPortsManager().PrepareConnection(net.ParseIP(c.ClientIP()))
	if !got { //没能预留连接。 If haven't reserved connection
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
	/*get parameters*/
	offset, err := strconv.ParseUint(c.GetHeader("offset"), 10, 64)
	filename := c.GetHeader("filename")
	clientFilePath := c.GetHeader("clientFilePath")
	path := c.GetHeader("path")
	uploadPath := path + filename
	fileSize, err := strconv.ParseUint(c.GetHeader("fileSize"), 10, 64)
	uploadId, err := strconv.Atoi(c.GetHeader("uploadId"))
	userId := GetUserIdFromContext(c)
	if len(path) == 0 || len(filename) == 0 || len(clientFilePath) == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing required info \n" + err.Error(),
		})
		return
	}
	err = Service.Upload(c, offset, uploadPath, fileSize, uploadId, userId, clientFilePath)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "uploaded",
		})
	}
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
	connection.GetDS2CS() <- -1
	///*设置响应头*/
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Connection", "close")
	///*写入输出流*/
	//chaDecipher := Utils.DefaultChaEncryptor()
	/*获取文件reader*/
	file, err := os.Open(Service.GetFullFilePath(filePath, userId))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	cha20FileIO, err := Utils.DefaultChaCha20FileIO(file, nil)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//Set the offset to 24, because the nonce put in the head of the encrypt file is 24-Bytes long.
	var offset int64 = 24
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
				offset = opcode + 24
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
			n, err := cha20FileIO.ReadAt(plaintext, file, offset)
			if err != nil && err.Error() != "EOF" { //文件读取出错
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": err.Error(),
				})
				break
			}
			if n == 0 {
				break //跳出for循环
			}
			_, err = w.Write(plaintext[:n])
			if err != nil {
				return false
			}
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

/*upload encrypt*/
//plaintext := make([]byte, buf)
//var receivedBytes uint64 = 0
//var cipherBuffer bytes.Buffer
//plaintFile, _ := os.Open("/Users/hanminghao/Desktop/study/test.txt")
//for {
//n, readErr := plaintFile.Read(plaintext)
//if readErr != nil && readErr != io.EOF && !errors.Is(readErr, http.ErrBodyReadAfterClose) { //读取出错
//break
//}
//if n == 0 {
//break
//}
//loopTimes := int((int64(n) + sectionSize - 1) / sectionSize)
//wg := sync.WaitGroup{}
//wg.Add(loopTimes)
//for index := 0; index < loopTimes; index++ {
//go func(idx int64) {
//upto := min((idx+1)*sectionSize, int64(n))
//Utils.RabbitEncrypt(plaintext[idx*sectionSize:upto], ciphertext[idx*sectionSize:upto])
//wg.Done()
//}(int64(index))
//}
//wg.Wait()
//if _, err := file.Write(ciphertext[:n]); err != nil {
//break
//}
//receivedBytes += uint64(n)
//}

/**/
//向上取整
//loopTimes := int((int64(n) + sectionSize - 1) / sectionSize)
//wg := sync.WaitGroup{}
//wg.Add(loopTimes)
//for index := 0; index < loopTimes; index++ {
//go func(idx int64) {
//upto := min((idx+1)*sectionSize, int64(n))
////chaDecipher.Decrypt(ciphertext[idx*sectionSize:upto], plaintext[idx*sectionSize:upto])
//Utils.RabbitDecrypt(plaintext[idx*sectionSize:upto], ciphertext[idx*sectionSize:upto])
////copy(plaintext[index*sectionSize:upto], ciphertext[index*sectionSize:upto])
//wg.Done()
//}(int64(index))
//}
//wg.Wait()
/**/
