package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/admin"
	"uapply_go/web/api/super_admin"
	"uapply_go/web/middleware"
)

// InitAdminRouter 初始化具体的路由组
func InitAdminRouter(router *gin.RouterGroup) {
	uGroup := router.Group("/admin")
	AdminWithJwt := router.Group("/admin").Use(middleware.JWTAuth())
	SuperAdmin := router.Group("/admin").Use(middleware.JWTAuth(), middleware.SupAdmin())
	{
		// Post 和 Put 接受表单参数
		uGroup.POST("/superadmin", super_admin.Create) // 超级管理员（组织）的创建,该接口不对外暴露，由后端人员操作
		SuperAdmin.POST("/admin", admin.Create)        // 部门（admin）的创建
		uGroup.POST("/login", admin.Login)             // 部门或组织登录

		AdminWithJwt.PUT("/org/dep", admin.Update) // 根据部门id更新部门信息
		SuperAdmin.PUT("/org", super_admin.Update) // 根据组织id更新组织信息

		// Get和Delete接收path参数
		uGroup.GET("/org/:orgid", super_admin.GetOrg)               // 通过id获取某一组织的信息,这个不需要权限，想获取就获取
		uGroup.GET("/:orgid", super_admin.GetOrgDep)                // 通过id获取某一组织下的所有部门
		uGroup.GET("/dep/:depid", admin.GetDetail)                  // 获取某一个部门的详细信息
		SuperAdmin.GET("/org/dep/:depid", super_admin.GetDepDetail) // 最高权限获取部门信息

		//uGroup.DELETE("/superadmin") //  超级管理员（组织）的注销,暂时不做支持
		SuperAdmin.DELETE("/admin/:depid", super_admin.Delete) // admin，即部门的删除
	}
}
