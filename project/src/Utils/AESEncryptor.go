package Utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type AesEncryptor struct {
	key []byte
	iv  []byte
}

var aesEncryptor *AesEncryptor

func DefaultAESEncryptor() *AesEncryptor {
	if aesEncryptor == nil {
		aesEncryptor = &AesEncryptor{}
	}
	keyIvYmlPath := DefaultConfigReader().Get("keyIvYmlPath").(string)
	keyStr := YmlReader(keyIvYmlPath, "key").(string)
	ivStr := YmlReader(keyIvYmlPath, "iv").(string)
	aesEncryptor.key = []byte(keyStr)
	aesEncryptor.iv = []byte(ivStr)
	return aesEncryptor
}

func (aesEncryptor *AesEncryptor) Encrypt(plaintext []byte, ciphertext []byte) error {
	if aesEncryptor.key == nil || aesEncryptor.iv == nil {
		return fmt.Errorf("no key and iv")
	}
	block, err := aes.NewCipher(aesEncryptor.key)
	if err != nil {
		return err
	}
	mode := cipher.NewCBCEncrypter(block, aesEncryptor.iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return nil
}

func (aesEncryptor *AesEncryptor) Decrypt(ciphertext []byte, plaintext []byte) error {
	if aesEncryptor.key == nil || aesEncryptor.iv == nil {
		return fmt.Errorf("no key and iv")
	}
	block, err := aes.NewCipher(aesEncryptor.key)
	if err != nil {
		return err
	}
	mode := cipher.NewCBCDecrypter(block, aesEncryptor.iv)
	
	mode.CryptBlocks(plaintext, ciphertext)
	return nil
}

func Encrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext = pkcs7Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func Decrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = pkcs7UnPadding(plaintext)
	return plaintext, nil
}

// PKCS#7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS#7 去填充
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

/**
使用示例
*/
//func main() {
//	key := []byte("this is a 16 byte key")
//	iv := []byte("this is a 16 byte iv")
//
//	plaintext := []byte("hello world")
//
//	// 加密
//	ciphertext, err := encrypt(plaintext, key, iv)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("加密结果：%s\n", base64.StdEncoding.EncodeToString(ciphertext))
//
//	// 解密
//	decrypted, err := decrypt(ciphertext, key, iv)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("解密结果：%s\n", decrypted)
//}
