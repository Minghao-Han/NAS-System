package Utils

import (
	"crypto/cipher"
	"fmt"
	"snix.ir/rabbitio"
)

var (
	key = []byte("key-gen-rabbitio")
	ivx = []byte("abcd8795")
)

func DefaultRabbitCipher() (cipher.Stream, error) {
	str, err := rabbitio.NewCipher(key, ivx)
	if err != nil {
		return nil, err
	}
	return str, nil
}
func RabbitEncrypt(plaintext []byte, ciphertext []byte) {
	str, err := rabbitio.NewCipher(key, ivx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			recover()
			fmt.Println("failed to encrypt")
		}
	}()
	str.XORKeyStream(ciphertext, plaintext)
}

func RabbitDecrypt(plaintext []byte, ciphertext []byte) {
	str, err := rabbitio.NewCipher(key, ivx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			recover()
			fmt.Println("failed to decrypt")
		}
	}()
	str.XORKeyStream(plaintext, ciphertext)
}

//
//func main() {
//	ptx := "plain text -- dummy text to encrypt and decrypt with rabbit"
//	str, err := rabbitio.NewCipher(key, ivx)
//	if err != nil {
//		panic(err)
//	}
//
//	cpt := make([]byte, len(ptx))
//	str.XORKeyStream(cpt, []byte(ptx))
//	//fmt.Println("cipher text ---:", hex.EncodeToString(cpt))
//
//	str, err = rabbitio.NewCipher(key, ivx)
//	if err != nil {
//		panic(err)
//	}
//
//	// decrypt cipher text and print orginal text
//	plx := make([]byte, len(cpt))
//	str.XORKeyStream(plx, cpt)
//	fmt.Println("plain text ----:", string(plx))
//}
