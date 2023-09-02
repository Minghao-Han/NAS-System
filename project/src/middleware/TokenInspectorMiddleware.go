package middleware

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	"net/http"
)

func TokenInspect() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "no token",
			})
			c.Abort()
			return
		}
		tokenValid, _ := Utils.DefaultJWT().VerifyToken(token)
		if !tokenValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "token invalid",
			})
			c.Abort()
			return
		}
		value, err := Utils.DefaultJWT().ExtractClaim(token, "userId")
		userId := convertToCorrectType(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "token doesn't contain user id",
			})
			c.Abort()
			return
		}
		c.Set("userId", userId)
	}
}

func convertToCorrectType(value interface{}) interface{} {
	switch v := value.(type) {
	case float64:
		if v == float64(int(v)) {
			return int(v)
		}
	}
	return value
}
