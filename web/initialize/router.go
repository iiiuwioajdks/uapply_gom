package initialize

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/middleware"
	"uapply_go/web/router"
)

// InitRouter 路由初始化
func InitRouter() *gin.Engine {
	Router := gin.New()
	Router.Use(GinLogger(), GinRecovery(true), middleware.Cors())

	ApiRouter := Router.Group("/uapply")
	router.InitUserRouter(ApiRouter)
	router.InitAdminRouter(ApiRouter)
	router.InitSuperAdminRouter(ApiRouter)
	router.InitCommonRouter(ApiRouter)
	return Router
}
