package initialize

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/router"
)

// InitRouter 路由初始化
func InitRouter() *gin.Engine {
	Router := gin.New()
	Router.Use(GinLogger(), GinRecovery(true))

	ApiRouter := Router.Group("/uapply/v1")
	router.InitUserRouter(ApiRouter)

	return Router
}
