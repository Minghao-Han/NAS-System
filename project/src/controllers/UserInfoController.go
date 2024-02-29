package controllers

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Service"
	"nas/project/src/Utils"
	"net/http"
)

func Test(c *gin.Context) {
	//userId := 0
	//user, _ := Service.GetUser(userId)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg":      "successfully get user info",
	//	"userInfo": user,
	//})
	c.JSON(http.StatusOK, gin.H{
		"msg": "The back-end is running now!",
	})
	return
}

type loginData struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func Login(c *gin.Context) {
	ld := loginData{}
	if err := c.ShouldBindJSON(&ld); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	userId, verified := Service.Authenticate(ld.Username, ld.Password)
	if !verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "wrong Username or wrong Password",
		})
	} else {
		token, err := Utils.DefaultJWT().GenerateToken(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "error generating token",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":   "successfully login",
				"token": token,
			})
		}
	}
}

func GetUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, err := Service.GetUser(userId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "user doesn't exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "successfully get user info",
		"userInfo": user,
	})
	return
}
