package Service

import (
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
		return -1, false
	}
	aesEncryptor := Utils.DefaultAESEncryptor()
	decryptPw, err := aesEncryptor.DecryptWithUnpadding(gottenPw) //解码
	if err != nil {
		return -1, false
	}
	if string(decryptPw) == password {
		user, err := userDA.FindByUsername(username) //获取密码
		if err != nil || user == nil {
			return -1, false
		}
		return user.UserId, true
	}
	return -1, false
}
