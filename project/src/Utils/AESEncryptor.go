package Utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AesEncryptor struct {
	key []byte
	iv  []byte
}

var defaultAesEncryptor *AesEncryptor

func DefaultAESEncryptor() *AesEncryptor {
	if defaultAesEncryptor == nil {
		defaultAesEncryptor = &AesEncryptor{}
		keyIvYmlPath := DefaultConfigReader().Get("Aes:keyIvYmlPath").(string)
		keyStr := YmlReader(keyIvYmlPath, "key").(string)
		ivStr := YmlReader(keyIvYmlPath, "iv").(string)
		defaultAesEncryptor.key = []byte(keyStr)
		defaultAesEncryptor.iv = []byte(ivStr)
	}
	return defaultAesEncryptor
}

func (aesEncryptor *AesEncryptor) EncryptWithPadding(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesEncryptor.key)
	if err != nil {
		return nil, err
	}
	plaintext = pkcs7Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, aesEncryptor.iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func (aesEncryptor *AesEncryptor) DecryptWithUnpadding(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesEncryptor.key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, aesEncryptor.iv)
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = pkcs7UnPadding(plaintext)
	return plaintext, nil
}

func AesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
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

func AesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
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
//	Nonce := []byte("this is a 16 byte Nonce")
//
//	plaintext := []byte("hello world")
//
//	// 加密
//	ciphertext, err := encrypt(plaintext, key, Nonce)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("加密结果：%s\n", base64.StdEncoding.EncodeToString(ciphertext))
//
//	// 解密
//	decrypted, err := decrypt(ciphertext, key, Nonce)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("解密结果：%s\n", decrypted)
//}
