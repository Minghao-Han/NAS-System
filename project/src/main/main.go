package main

import (
	"nas/project/src/Utils"
	"nas/project/src/router"
	"sync"
)

func main() {
	normalRouter := router.GetNormalRouter()
	csPorts := Utils.DefaultConfigReader().Get("FSP:csPorts").([]interface{})
	dsPorts := Utils.DefaultConfigReader().Get("FSP:dsPorts").([]interface{})
	wg := sync.WaitGroup{}
	wg.Add(len(csPorts) + len(dsPorts) + 1)
	for _, csPort := range csPorts {
		csPort := csPort
		csRouter := router.GetControlStreamRouter()
		go func() {
			router.RunOnConfig(csPort.(int), csRouter)
			wg.Done()
		}()
	}
	for _, dsPort := range dsPorts {
		dsPort := dsPort
		dsRouter := router.GetDataStreamRouter()
		go func() {
			router.RunOnConfig(dsPort.(int), dsRouter)
			wg.Done()
		}()
	}
	go func() {
		router.RunTLSOnConfig(normalRouter)
		wg.Done()
	}()
	wg.Wait()
}

/*---------------------------------------------------------------*/
//func main() {
//	ciphertextBuffer := encryptFile()
//	cipherFile, _ := os.Create("/Users/hanminghao/Desktop/study/standard_ciphertext.txt")
//	defer cipherFile.Close()
//	writer := bufio.NewWriter(cipherFile)
//	ciphertextBuffer.WriteTo(writer)
//	writer.Flush()
//	plaintext := decryptFile(ciphertextBuffer)
//	fmt.Println(plaintext.String())
//}
//func encryptFile() (aa bytes.Buffer) {
//	sectionSize := 5000 / int64(10)
//	plaintext := make([]byte, 9000)
//	ciphertext := make([]byte, 9000)
//	var receivedBytes uint64 = 0
//	file, _ := os.Open("/Users/hanminghao/Desktop/study/test.txt")
//	for {
//		n, readErr := file.Read(plaintext)
//		if readErr != nil && readErr != io.EOF && !errors.Is(readErr, http.ErrBodyReadAfterClose) { //读取出错
//			break
//		}
//		if n == 0 {
//			break
//		}
//		loopTimes := int((int64(n) + sectionSize - 1) / sectionSize)
//		wg := sync.WaitGroup{}
//		wg.Add(loopTimes)
//		for index := 0; index < loopTimes; index++ {
//			go func(idx int64) {
//				upto := min((idx+1)*sectionSize, int64(n))
//				Utils.RabbitEncrypt(plaintext[idx*sectionSize:upto], ciphertext[idx*sectionSize:upto])
//				wg.Done()
//			}(int64(index))
//		}
//		wg.Wait()
//		if _, err := aa.Write(ciphertext[:n]); err != nil {
//			break
//		}
//		receivedBytes += uint64(n)
//	}
//	return aa
//}
//func decryptFile(buffer bytes.Buffer) (result bytes.Buffer) {
//	cipherFile, _ := os.Open("/Users/hanminghao/Desktop/study/standard_ciphertext.txt")
//	sectionSize := 5000 / int64(10)
//	var offset int64 = 0
//	ciphertext := make([]byte, 5000)
//	plaintext := make([]byte, 5000)
//	for {
//		/*解密*/
//		n, err := cipherFile.Read(ciphertext)
//		if err != nil && err.Error() != "EOF" { //文件读取出错
//			break
//		}
//		if n == 0 {
//			break //跳出for循环
//		}
//		/**/
//		//向上取整
//		loopTimes := int((int64(n) + sectionSize - 1) / sectionSize)
//		wg := sync.WaitGroup{}
//		wg.Add(loopTimes)
//		fmt.Println("download")
//		for index := 0; index < loopTimes; index++ {
//			go func(idx int64) {
//				upto := min((idx+1)*sectionSize, int64(n))
//				//chaDecipher.Decrypt(ciphertext[idx*sectionSize:upto], plaintext[idx*sectionSize:upto])
//				Utils.RabbitDecrypt(plaintext[idx*sectionSize:upto], ciphertext[idx*sectionSize:upto])
//				wg.Done()
//			}(int64(index))
//		}
//		wg.Wait()
//		/**/
//		result.Write(plaintext[:n])
//		offset += int64(n)
//	}
//	return result
//}

/*---------------------------------------------------------------*/

//func compareFiles(file1Path, file2Path string) (bool, error) {
//	// 打开第一个文件
//	result := true
//	file1, err := os.Open(file1Path)
//	if err != nil {
//		return false, fmt.Errorf("无法打开文件1: %v", err)
//	}
//	defer file1.Close()
//
//	// 打开第二个文件
//	file2, err := os.Open(file2Path)
//	if err != nil {
//		return false, fmt.Errorf("无法打开文件2: %v", err)
//	}
//	defer file2.Close()
//
//	// 创建带缓冲的读取器
//	reader1 := bufio.NewReader(file1)
//	reader2 := bufio.NewReader(file2)
//
//	// 逐行比较文件内容
//	lineNum := 1
//	for {
//		line1, err1 := reader1.ReadString('\n')
//		line2, err2 := reader2.ReadString('\n')
//
//		if line1 != line2 {
//			// 发现不同的行
//			fmt.Printf("第%d行不同\n", lineNum)
//			return false, nil
//		}
//
//		if err1 != nil || err2 != nil {
//			// 检查两个文件是否都已到达末尾
//			if err1.Error() == "EOF" && err2.Error() == "EOF" {
//				break
//			} else {
//				return false, fmt.Errorf("读取文件时发生错误: 文件1(%v), 文件2(%v)", err1, err2)
//			}
//		}
//
//		lineNum++
//	}
//
//	return result, nil
//}
//
//func main() {
//	file1Path := "/Users/hanminghao/Desktop/study/standard_ciphertext.txt"
//	file2Path := "/Users/hanminghao/Desktop/study/NAS_Disk_Root/0/new/test_2.txt"
//
//	same, err := compareFiles(file1Path, file2Path)
//	if err != nil {
//		fmt.Println("比较文件时发生错误:", err)
//		return
//	}
//
//	if same {
//		fmt.Println("两个文件的内容相同.")
//	} else {
//		fmt.Println("两个文件的内容不相同.")
//	}
//}
