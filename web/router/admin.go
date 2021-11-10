package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/admin"
	"uapply_go/web/middleware"
)

// InitAdminRouter 初始化具体的路由组
// 后台管理系统相关接口
func InitAdminRouter(router *gin.RouterGroup) {
	adminr := router.Group("/dep").Use(middleware.JWTAuth())
	{
		// 首先需要判断 num 的值，只能为1和2，1是第一轮面试，2 是第二轮面试
		adminr.POST("/interviewer/add")       // 增加面试官
		adminr.POST("/interview/sms/:num")    // 发送第n轮面试短信
		adminr.POST("/interview/email/:num")  // 发送第n轮面试邮件
		adminr.POST("/interview/pass/:num")   // 通过第n轮面试
		adminr.POST("/interview/out/:num")    // 在第n轮面试时直接淘汰，进行数据库软删除
		adminr.POST("/interview/enroll/:num") // 在第n轮面试时加入暂录取名单，加入部员名单

		adminr.PATCH("/udpdep", admin.Update) // 根据部门id更新部门信息

		adminr.GET("/getdep", admin.Get)
		adminr.GET("/getusers/register")     // 部门获取报名自己部门的所有用户
		adminr.GET("/getuser/register")      // 部门获取报名自己部门的某一个用户详细信息
		adminr.GET("/getuser/unreview/:num") // 部门获取第n轮未面试成员
		adminr.GET("/getuser/reviewed/:num") // 部门获取第n轮已面试成员
		adminr.GET("/getuser/enroll")        // 部门获取自己的通过部员
		adminr.GET("/getuser/info")          // 获取本部门男女人数，报名人数信息

		adminr.DELETE("/interviewer/del") // 删除面试官
	}
}
