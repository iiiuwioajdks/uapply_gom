package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/user"
	"uapply_go/web/middleware"
)

// InitUserRouter 初始化具体的路由组
func InitUserRouter(router *gin.RouterGroup) {
	uGroup := router.Group("/user")

	{
		// code 是前端给我们，是微信登录凭证
		uGroup.GET("/login/:code", user.Login)
	}

	jwtUser := router.Group("/user", middleware.WXJWTAuth())

	{
		// 下面这些接口如果需要操作数据库，就先 c.get("wxClaim") 拿到 openid,然后去查找 uid，再去得到对应的信息
		// 要看详细去阅读 middleware jwt 中的 WXJWTAuth 源码
		jwtUser.POST("/resume/save", user.SaveResume)       // 用户保存自己的简历
		jwtUser.POST("/resume/tmpsave", user.SaveTmpResume) // 用户保存简历到草稿箱
		jwtUser.POST("/register", user.Register)            // 用户报名

		jwtUser.GET("/register/info", user.GetRegInfo)     // 用户查看已经提交的报名信息，及用户报名的组织下的部门
		jwtUser.GET("/register/status", user.GetRegStatus) // 用户查看已经报名的部门的招新状态
		jwtUser.GET("/resume/gettmp", user.GetTmpResume)   // 用户获取草稿箱简历
		jwtUser.GET("/resume/get", user.GetResume)         // 用户获取已经填写的简历信息

		jwtUser.PATCH("/resume", user.UpdateResume)   // 用户更新简历部分内容
		jwtUser.PATCH("/resume/text", user.ClearText) // 一建清除text文本，这个要在数据库直接插入空，单独开一个接口可以省去很多判断
	}
}
