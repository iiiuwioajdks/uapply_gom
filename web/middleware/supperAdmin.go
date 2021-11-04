package middleware

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api"
	j "uapply_go/web/models/jwt"
)

func SupAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, exist := c.Get("claim")
		if !exist {
			api.Fail(c, api.CodeUserInfoNotExist)
			c.Abort()
			return
		}
		claimInfo := claim.(*j.Claims)
		if claimInfo.Role != 1 {
			api.Fail(c, api.CodeHasNotPower)
			c.Abort()
			return
		}
		c.Next()
	}
}
