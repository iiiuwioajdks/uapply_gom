package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/interviewer"
	"uapply_go/web/middleware"
)

// InitInterviewerRouter 小程序端面试官相关路由
func InitInterviewerRouter(router *gin.RouterGroup) {
	iRouter := router.Group("/interviewer")
	{
		iRouter.GET("/login/:code", interviewer.Login)
	}

	jwtIRouter := router.Group("/interviewer").Use(middleware.WXJWTAuth(), middleware.IsInterviewer())
	{
		// 这里前端就要负责把 Position 中拿到的，然后用户选进去的部门的id告诉我了，这样我才能提交
		jwtIRouter.POST("/evaluate", interviewer.Evaluate) // 对用户的评价提交

		// 获得面试官所在的组织和部门，一个微信（人）可能是多个组织的面试官,通过jwt拿到id就行，不用前端传
		// 这里要负责把depid和orgid传回给前端
		jwtIRouter.GET("/position", interviewer.Position)
		jwtIRouter.GET("/getuser/:uid", interviewer.GetUser)   // 获取用户简历,这个uid是用户的uid
		jwtIRouter.GET("/info/:depid", interviewer.GetDepInfo) // 获取部门投递人数，性别，通过人数等信息，首先要判断这个部门对不对应这个人，这个人能不能搞这个部门
	}
}
