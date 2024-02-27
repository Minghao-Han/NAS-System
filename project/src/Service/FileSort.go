package Service

import (
	"os"
	"sort"
)

/*
修改时间排序
*/
func FileInfosByModifiedTimeUp(fileInfos []os.FileInfo) []os.FileInfo {
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime().After(fileInfos[j].ModTime())
	})
	return fileInfos
}
func FileInfosByModifiedTimeDown(fileInfos []os.FileInfo) []os.FileInfo {
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime().Before(fileInfos[j].ModTime())
	})
	return fileInfos
}

/*文件大小排序*/
func FileInfosBySizeUp(fileInfos []os.FileInfo) []os.FileInfo {
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Size() > fileInfos[j].Size()
	})
	return fileInfos
}
func FileInfosBySizeDown(fileInfos []os.FileInfo) []os.FileInfo {
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Size() < fileInfos[j].Size()
	})
	return fileInfos
}

/*文件字典排序*/
func FileInfosByDicUp(fileInfos []os.FileInfo) []os.FileInfo {
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})
	return fileInfos
}
func FileInfosByDicDown(fileInfos []os.FileInfo) []os.FileInfo {
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() > fileInfos[j].Name()
	})
	return fileInfos
}
