package router

import "github.com/gin-gonic/gin"

// InitUserRouter 初始化具体的路由组
func InitUserRouter(router *gin.RouterGroup) {
	uGroup := router.Group("/user")

	{
		uGroup.GET("/")
	}
}
