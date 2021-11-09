package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/admin"
	"uapply_go/web/middleware"
)

// InitAdminRouter 初始化具体的路由组
func InitAdminRouter(router *gin.RouterGroup) {
	adminr := router.Group("/dep").Use(middleware.JWTAuth())
	{
		adminr.PATCH("/udpdep", admin.Update) // 根据部门id更新部门信息
		adminr.GET("/getdep", admin.Get)
	}
}
