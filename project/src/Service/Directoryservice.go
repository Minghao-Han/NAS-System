package Service

import (
	"fmt"
	"nas/project/src/Utils"
	"os"
	"strconv"
	"time"
)

type FileInfoVO struct {
	Name    string    `json:"name,omitempty"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

var pageSize = Utils.DefaultConfigReader().Get("CheckDir:pageSize").(int)

func CheckDir(userId int, dirPath string, order string, pageNum int, catalog string) ([]FileInfoVO, error) {
	dirFullPath := diskRoot + strconv.Itoa(userId) + dirPath
	dir, err := os.Open(dirFullPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fileInfos, err := ReadDir(dir, catalog)
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
func toFileInfo(osFileInfos []os.FileInfo) []FileInfoVO {
	toBeReturn := make([]FileInfoVO, 0, len(osFileInfos))
	for _, osFileInfo := range osFileInfos {
		toBeReturn = append(toBeReturn, FileInfoVO{
			Name:    osFileInfo.Name(),
			Size:    osFileInfo.Size(),
			ModTime: osFileInfo.ModTime(),
		})
	}
	return toBeReturn
}
func CreateDir(parentDirPath string, dirName string, userId int) error {
	parentDirPath = GetFullFilePath(parentDirPath, userId)
	_, err := os.Stat(parentDirPath)
	if os.IsNotExist(err) {
		// 文件夹不存在
		return fmt.Errorf("parent folder does not exist")
	}
	dirPath := parentDirPath + Slash + dirName
	err = os.Mkdir(dirPath, 0700)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDir(dirPath string, userId int) error {
	dirPath = GetFullFilePath(dirPath, userId)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist")
	}
	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}
	return nil
}

func ReadDir(dir *os.File, catalog string) ([]os.FileInfo, error) {
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	switch catalog {
	case "all":
		return fileInfos, nil
	case "image":
		result := make([]os.FileInfo, 0)
		for _, fileInfo := range fileInfos {
			if Utils.IsImage(fileInfo.Name()) {
				result = append(result, fileInfo)
			}
		}
		return result, nil
	case "video":
		result := make([]os.FileInfo, 0)
		for _, fileinfo := range fileInfos {
			if Utils.IsVideo(fileinfo.Name()) {
				result = append(result, fileinfo)
			}
		}
		return result, nil
	case "doc":
		result := make([]os.FileInfo, 0)
		for _, fileinfo := range fileInfos {
			if Utils.IsDoc(fileinfo.Name()) {
				result = append(result, fileinfo)
			}
		}
		return result, nil
	case "other":
		result := make([]os.FileInfo, 0)
		for _, fileinfo := range fileInfos {
			if !Utils.IsDoc(fileinfo.Name()) || !Utils.IsImage(fileinfo.Name()) || !Utils.IsVideo(fileinfo.Name()) {
				result = append(result, fileinfo)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("catalog invalid")
	}
}
