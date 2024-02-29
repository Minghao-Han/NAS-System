package Utils

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"io"
	"os"
)

type Cha20IO struct {
	key       []byte
	encoder   *chacha20.Cipher
	decoder   *chacha20.Cipher
	bufReader *bufio.Reader
	bufWriter *bufio.Writer
}

func DefaultChaCha20FileIO(reader io.Reader, writer io.Writer) (*Cha20IO, error) {
	keyYmlPath := DefaultConfigReader().Get("ChaCha20:keyNonceYmlPath").(string)
	keyStr := YmlReader(keyYmlPath, "key").(string)
	s := &Cha20IO{
		key: []byte(keyStr), // should be exactly 32 bytes
	}

	var err error
	nonce := make([]byte, chacha20.NonceSizeX)
	//要区分是新创建的空文件还是老文件再追加。如果是新，rand一个nonce；如果是旧，读之前的nonce
	n, _ := reader.Read(nonce)
	if n == 0 { //indicate that this is a new file, we need to generate a random nonce for it.
		if _, err := rand.Read(nonce); err != nil {
			return nil, err
		}
		s.bufWriter = bufio.NewWriter(writer)
		if n, err := s.bufWriter.Write(nonce); err != nil || n != len(nonce) {
			return nil, errors.New("write nonce failed: " + err.Error())
		}
		err := s.bufWriter.Flush()
		if err != nil {
			return nil, err
		}
	} else if n != chacha20.NonceSizeX {
		return nil, fmt.Errorf("file corruption")
	}

	s.encoder, err = chacha20.NewUnauthenticatedCipher(s.key, nonce)
	if err != nil {
		return nil, err
	}

	return s, nil
}
func NewCha20IO(key []byte, reader io.Reader, writer io.Writer) (*Cha20IO, error) {
	s := &Cha20IO{
		key: key, // should be exactly 32 bytes
	}

	var err error
	nonce := make([]byte, chacha20.NonceSizeX)
	//要区分是新创建的空文件还是老文件再追加。如果是新，rand一个nonce；如果是旧，读之前的nonce
	n, _ := reader.Read(nonce)
	if n == 0 { //indicate that this is a new file, we need to generate a random nonce for it.
		if _, err := rand.Read(nonce); err != nil {
			return nil, err
		}
		s.bufWriter = bufio.NewWriter(writer)
		if n, err := s.bufWriter.Write(nonce); err != nil || n != len(nonce) {
			return nil, errors.New("write nonce failed: " + err.Error())
		}
		err := s.bufWriter.Flush()
		if err != nil {
			return nil, err
		}
	} else if n != chacha20.NonceSizeX {
		return nil, fmt.Errorf("file corruption")
	}

	s.encoder, err = chacha20.NewUnauthenticatedCipher(s.key, nonce)
	if err != nil {
		return nil, err
	}

	return s, nil
}
func (s *Cha20IO) Read(p []byte, reader io.Reader) (int, error) {
	s.bufReader = bufio.NewReader(reader)
	n, err := s.bufReader.Read(p)
	//n, err := s.buffer.Read(p)
	if err != nil || n == 0 {
		return n, err
	}

	dst := make([]byte, n)
	s.decoder.XORKeyStream(dst, p[:n])
	copy(p[:n], dst)
	return n, nil
}
func (s *Cha20IO) ReadAt(p []byte, file *os.File, offset int64) (int, error) { //Only for file io
	if s.decoder == nil {
		nonce := make([]byte, chacha20.NonceSizeX)
		//if n, err := s.buffer.Read(nonce); err != nil || n != len(nonce) {
		if n, err := file.ReadAt(nonce, 0); err != nil || n != len(nonce) {
			return n, errors.New("can't read nonce from stream: " + err.Error())
		}
		decoder, err := chacha20.NewUnauthenticatedCipher(s.key, nonce)
		if err != nil {
			return 0, errors.New("generate decoder failed: " + err.Error())
		}
		s.decoder = decoder
	}
	_, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return s.Read(p, file)
}
func (s *Cha20IO) Write(p []byte, writer io.Writer) (int, error) {
	dst := make([]byte, len(p))
	s.encoder.XORKeyStream(dst, p)
	if s.bufWriter == nil {
		s.bufWriter = bufio.NewWriter(writer)
	}
	return s.bufWriter.Write(dst)
}

func (s *Cha20IO) Close() error {
	if s.bufWriter != nil {
		return s.bufWriter.Flush()
	}
	s.bufReader = nil
	s.bufWriter = nil
	s.encoder = nil
	s.decoder = nil
	return nil
}

/*example of encryption and decryption*/
func encrypt() {
	/*1. Get the input io.Reader*/
	filePath := "plain_text_file_path"
	file, err := os.Open(filePath)
	plainTextReader := io.Reader(file)
	/*2. Create or open the encrypt file*/
	encryptFile, err := os.Create("encrypt_text_file_path")
	defer func() {
		file.Close()
		encryptFile.Close()
	}()
	/*3. new an instance of Cha20IO*/
	cha20IO, err := NewCha20IO([]byte("a1b2c3d4e5f6g7h8i9j0k1l2m3n4o522"), encryptFile, encryptFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	/*4. Make a []byte to store read plain text. The size of it depends on the capacity of your memory*/
	plaintext := make([]byte, 10)
	/*5. Read plain text in the loop*/
	for {
		//read
		n, _ := plainTextReader.Read(plaintext)
		if n == 0 {
			break
		}
		/*6. Write the plain text to file. The cha20IO.Write() will automatically encrypt it and then write to the file*/
		_, err := cha20IO.Write(plaintext[:n], encryptFile)
		if err != nil {
			panic(err)
		}
	}
	/*7. Close the cha20IO*/
	cha20IO.Close()
}
func decrypt() {
	/*1. Open the encrypt file*/
	encryptFile, _ := os.Open("encrypt_text_file_path")
	/*2. Get the output io.Writer*/
	decryptFile, _ := os.Create("decrypt_text_file_path")
	decryptWriter := bufio.NewWriter(decryptFile)
	defer func() {
		//decryptWriter.Flush()
		encryptFile.Close()
		decryptFile.Close()
	}()
	/*3. new an instance of Cha20IO*/
	cha20IO, err := NewCha20IO([]byte("a1b2c3d4e5f6g7h8i9j0k1l2m3n4o522"), encryptFile, encryptFile)
	if err != nil {
		panic(err)
	}
	/*4. Make a []byte to store read plain text. The size of it depends on the capacity of your memory*/
	p := make([]byte, 8)
	/*5. Set the offset to 24, because the nonce put in the head of the encrypt file is 24-Bytes long. The offset decide where you want to read
	Don't forget to increase the offset in the loop!*/
	var offset int64 = 24
	for {
		/*6. Read from file. The cha20IO.ReadAt will automatically decrypt the file and fill the p parameter with plain text*/
		n, err := cha20IO.ReadAt(p, encryptFile, offset)
		if n == 0 || (err != nil && err != io.EOF) {
			//fmt.Println(err.Error())
			break
		}
		/*7. Write decrypt data to output*/
		_, err = decryptWriter.Write(p[:n])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		offset += int64(n)
	}
}
