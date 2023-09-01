package Service

import (
	"nas/project/src/DA/userDA"
	"nas/project/src/Entities"
)

func GetUser(userId int) (*Entities.User, error) {
	return userDA.FindById(userId)
}

func Authenticate(username string, password string) bool {
	
}
