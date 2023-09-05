package Utils

import (
	"fmt"
	"golang.org/x/crypto/chacha20"
)

type ChaChaEncryptor struct {
	key   []byte
	Nonce []byte
}

var defaultChaEncryptor *ChaChaEncryptor

func DefaultChaEncryptor() *ChaChaEncryptor {
	if defaultChaEncryptor == nil {
		keyIvYmlPath := DefaultConfigReader().Get("ChaCha20:keyNonceYmlPath").(string)
		keyStr := YmlReader(keyIvYmlPath, "key").(string)
		ivStr := YmlReader(keyIvYmlPath, "nonce").(string)
		defaultChaEncryptor = &ChaChaEncryptor{
			key:   []byte(keyStr),
			Nonce: []byte(ivStr),
		}
	}
	return defaultChaEncryptor
}

// Encrypt 使用 ChaCha20 对明文进行加密
func (cha *ChaChaEncryptor) Encrypt(plaintext []byte, ciphertext []byte) error {
	if len(plaintext) != len(ciphertext) {
		return fmt.Errorf("ciphertext length mismatch")
	}
	c, err := chacha20.NewUnauthenticatedCipher(cha.key, cha.Nonce)
	if err != nil {
		return err
	}
	c.XORKeyStream(ciphertext, plaintext)
	return nil
}

// Decrypt 使用 ChaCha20 解密密文
func (cha *ChaChaEncryptor) Decrypt(ciphertext []byte, plaintext []byte) ([]byte, error) {
	if len(plaintext) != len(ciphertext) {
		return nil, fmt.Errorf("plaintext length mismatch")
	}
	c, err := chacha20.NewUnauthenticatedCipher(cha.key, cha.Nonce)
	if err != nil {
		return nil, err
	}
	c.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
func ChaEncrypt(key, nonce, plaintext []byte) ([]byte, error) {
	c, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	c.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

// Decrypt 使用 ChaCha20 解密密文
func ChaDecrypt(key, nonce, ciphertext []byte) ([]byte, error) {
	c, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	c.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

//
//func main() {
//	key := make([]byte, 32)   // 密钥
//	nonce := make([]byte, 12) // 初始向量（nonce）
//	message := []byte("Hello, World!") // 明文
//
//	ciphertext, err := Encrypt(key, nonce, message)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	plaintext, err := Decrypt(key, nonce, ciphertext)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("Plaintext: %s", plaintext)
//}
