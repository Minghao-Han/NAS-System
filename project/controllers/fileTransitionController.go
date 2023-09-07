package controllers

import (
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

var portsManager = PortManage.DefaultPortsManager()

func CsForUpload(c *gin.Context) {
	//	/**/
	//	sourceIP := net.ParseIP(c.ClientIP())
	//	url := c.Request.URL
	//	_, portStr, _ := net.SplitHostPort(url.Host)
	//	port, err := strconv.Atoi(portStr)
	//	if err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"msg": err.Error(),
	//		})
	//		return
	//	}
	//	portEntity, found := portsManager.FindPort(port, 0)
	//	if !found {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"msg": "didn't find port info",
	//		})
	//		return
	//	}
	//	connection, reserved := portEntity.FindConnection(sourceIP)
	//	if !reserved {
	//		c.JSON(http.StatusUnprocessableEntity, gin.H{
	//			"msg": "reservation didn't find",
	//		})
	//		return
	//	}
	//	cs2ds := connection.GetCS2DS()
	//	/**/
}

func DsForUpload(c *gin.Context) {

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
	sourceIP := net.ParseIP(c.ClientIP())
	url := c.Request.URL
	_, portStr, _ := net.SplitHostPort(url.Host)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	portEntity, found := portsManager.FindPort(port, 0)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "didn't find port info",
		})
		return
	}
	connection, reserved := portEntity.FindConnection(sourceIP)
	if !reserved {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "reservation didn't find",
		})
		return
	}
	cs2ds := connection.GetCS2DS()
	cs2ds <- opcode
}

func DsForDownload(c *gin.Context) {
	/*验证connection预留信息*/
	sourceIP := net.ParseIP(c.ClientIP())
	_, portStr, _ := net.SplitHostPort(c.Request.Host)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	portEntity, found := portsManager.FindPort(0, port)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "didn't find port info",
		})
		return
	}
	connection, reserved := portEntity.FindConnection(sourceIP)
	if !reserved {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "reservation didn't find",
		})
		return
	}
	defer portEntity.DisConnectByIP(sourceIP) //断开连接
	connection.GetDS2CS() <- -1
	cs2ds := connection.GetCS2DS()
	/*获取文件reader*/
	filePath := c.GetHeader("filePath")
	value, _ := c.Get("userId")
	userId := value.(int)
	file, err := os.Open(Service.GetFullFilePath(filePath, userId))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer file.Close()
	/*设置响应头*/
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	/*写入输出流*/
	chaDecipher := Utils.DefaultChaEncryptor()
	var offset int64 = 0
	ciphertext := make([]byte, bufferSize)
	plaintext := make([]byte, bufferSize)
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
			return
		}
		/*解密*/
		n, err := file.ReadAt(ciphertext, offset)
		if err != nil && err.Error() != "EOF" { //文件读取出错
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		if n == 0 {
			break //跳出for循环
		}
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
		clientDisconnected := c.Stream(func(w io.Writer) bool {
			w.Write(plaintext[:n])
			return false
		})
		if clientDisconnected {
			break
		}
		offset += int64(n)
	}
	c.Stream(func(w io.Writer) bool { //结束文件流
		return true
	})
	return
}
