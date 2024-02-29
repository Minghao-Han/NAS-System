package Service

import (
	"fmt"
	"nas/project/src/DA/userDA"
	"nas/project/src/Entities"
	"nas/project/src/Utils"
)

func GetUser(userId int) (*Entities.User, error) {
	return userDA.FindById(userId)
}

func Authenticate(username string, password string) (int, bool) { //userId,密码正确否
	gottenPw, err := userDA.GetPassword(username)
	if err != nil {
		fmt.Println(err.Error())
		return -1, false
	}
	aesEncryptor := Utils.DefaultAESEncryptor()
	decryptPw, err := aesEncryptor.DecryptWithoutPadding(gottenPw) //解码
	if err != nil {
		fmt.Println(err.Error())
		return -1, false
	}
	if string(decryptPw) == password {
		user, err := userDA.FindByUsername(username) //获取密码
		if err != nil || user == nil {
			fmt.Println(err.Error())
			return -1, false
		}
		return user.UserId, true
	}
	return -1, false
}

func MarginAvailable(userId int, fileSize uint64) error {
	user, err := userDA.FindById(userId)
	if err != nil {
		return fmt.Errorf("user doesn't exist")
	}
	if user.Margin < fileSize {
		return fmt.Errorf("no more space for this file")
	}
	return nil
}
