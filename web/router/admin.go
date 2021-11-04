package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/super_admin"
)

// InitAdminRouter 初始化具体的路由组
func InitAdminRouter(router *gin.RouterGroup) {
	uGroup := router.Group("/admin")

	{
		// Post 和 Put 接受表单参数
		uGroup.POST("/superadmin", super_admin.Create) // 超级管理员（组织）的创建
		uGroup.POST("/admin")                          // 部门（admin）的创建
		uGroup.POST("/login")                          // 部门或组织登录

		uGroup.PUT("/org")     // 根据部门id更新部门信息
		uGroup.PUT("/org/dep") // 根据组织id更新组织信息

		// Get和Delete接收path参数

		uGroup.GET("/org/:orgid") // 通过id获取某一组织的信息
		uGroup.GET("/:orgid")     // 通过id获取某一组织下的所有部门
		uGroup.GET("/dep/:depid") // 获取某一个部门的详细信息

		//uGroup.DELETE("/superadmin") //  超级管理员（组织）的注销,暂时不做支持
		uGroup.DELETE("/admin/:depid") // admin，即部门的删除
	}
}
