package Service

import (
	"fmt"
	uploadDA "nas/project/src/DA/uploadLogDA"
	"nas/project/src/DA/userDA"
	"nas/project/src/Utils"
	"os"
	"strconv"
	"strings"
	"time"
)

var diskRoot = Utils.DefaultConfigReader().Get("DiskRoot").(string)

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
	newFilePath := destinyFullPath + "/" + fileInfo.Name()
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

type FileInfo struct {
	Name    string    `json:"name,omitempty"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

var pageSize = Utils.DefaultConfigReader().Get("CheckDir:pageSize").(int)

func CheckDir(userId int, dirPath string, order string, pageNum int) ([]FileInfo, error) {
	dirFullPath := diskRoot + strconv.Itoa(userId) + dirPath
	dir, err := os.Open(dirFullPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	switch order {
	case "earliest":
		fileInfos = FileInfosByModifiedTimeDown(fileInfos)
		break
	case "latest":
		fileInfos = FileInfosByModifiedTimeUp(fileInfos)
		break
	case "biggest":
		fileInfos = FileInfosBySizeUp(fileInfos)
		break
	case "smallest":
		fileInfos = FileInfosBySizeDown(fileInfos)
		break
	case "dic_down":
		fileInfos = FileInfosByDicDown(fileInfos)
		break
	case "dic_up":
		fileInfos = FileInfosByDicUp(fileInfos)
		break
	default:
		return nil, fmt.Errorf("order invalid")
	}
	from := pageSize * pageNum
	if from > len(fileInfos)-1 { //没有更多
		return nil, fmt.Errorf("no more files")
	}
	to := min(pageSize*(pageNum+1), len(fileInfos)) //slice左闭右开，所以是pageSize*(pageNum+1)而不是pageSize*(pageNum+1)-1
	fileInfos = fileInfos[from:to]
	return toFileInfo(fileInfos), nil
}
func toFileInfo(osFileInfos []os.FileInfo) []FileInfo {
	toBeReturn := make([]FileInfo, 0, len(osFileInfos))
	for _, osFileInfo := range osFileInfos {
		toBeReturn = append(toBeReturn, FileInfo{
			Name:    osFileInfo.Name(),
			Size:    osFileInfo.Size(),
			ModTime: osFileInfo.ModTime(),
		})
	}
	return toBeReturn
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
