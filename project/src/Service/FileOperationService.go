package Service

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/DA/userDA"
	"nas/project/src/Utils"
	"nas/project/src/Utils/ImageUtil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var diskRoot = Utils.DefaultConfigReader().Get("FileSystem:DiskRoot").(string)
var Slash = Utils.DefaultConfigReader().Get("FileSystem:SlashStyle").(string)

func DeleteFile(userId int, filePath string) error {
	fileFullPath := GetFullFilePath(filePath, userId)
	fileInfo, existErr := os.Stat(fileFullPath)
	if os.IsNotExist(existErr) { //文件不存在
		return existErr
	}
	rmErr := os.Remove(fileFullPath)
	if rmErr != nil { //删除文件失败
		return rmErr
	}
	/*还要删掉数据库中对应的行*/
	_, err := uploadDA.DeleteByPath(filePath) //用用户相对路径而不是跟路径，避免泄露
	if err != nil {
		return err
	}
	user, _ := userDA.FindById(userId)
	if user.Capacity-user.Margin < uint64(fileInfo.Size()) {
		return fmt.Errorf("delete err")
	}
	user.Margin += uint64(fileInfo.Size())
	userDA.Update(*user)
	return nil
}

func MoveFile(userId int, sourceFilePath string, destinyPath string, cover bool) error {
	userRoot := diskRoot + strconv.Itoa(userId)
	sourceFullPath := userRoot + sourceFilePath
	destinyFullPath := userRoot + destinyPath
	fileInfo, existErr := os.Stat(sourceFullPath)
	if os.IsNotExist(existErr) { //文件不存在
		return existErr
	}
	if _, dirErr := os.Stat(destinyFullPath); os.IsNotExist(dirErr) { //目的文件夹不存在
		return dirErr
	}
	newFilePath := destinyFullPath + Slash + fileInfo.Name()
	_, err := os.Stat(newFilePath)
	if !os.IsNotExist(err) { //有同名文件
		if !cover { //不覆盖
			return fmt.Errorf("duplicate file")
		}
	}
	renameErr := os.Rename(sourceFullPath, newFilePath)
	if renameErr != nil {
		return renameErr
	}
	return nil
}

func GetThumbnail(userId int, filePath string) (*bytes.Buffer, error) {
	filePath = GetFullFilePath(filePath, userId)
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	img, err := ImageUtil.ImgDecode(imgFile, filepath.Ext(filePath))
	if err != nil {
		return nil, err
	}
	thumbnail := resize.Resize(180, 0, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumbnail, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func FileExists(filePath string, userId int) bool {
	userRoot := diskRoot + strconv.Itoa(userId)
	fullFilePath := userRoot + filePath
	if _, err := os.Stat(fullFilePath); os.IsExist(err) { //文件存在
		return true
	}
	return false
}
func GetFullFilePath(filePath string, userId int) string {
	return diskRoot + strconv.Itoa(userId) + filePath
}
func GetUserRelativePath(fullFilePath string, userId int) string {
	return strings.Trim(fullFilePath, diskRoot+strconv.Itoa(userId))
}
