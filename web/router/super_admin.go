package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/admin"
	"uapply_go/web/api/super_admin"
	"uapply_go/web/middleware"
)

func InitSuperAdminRouter(router *gin.RouterGroup) {
	// token 认证，组织认证
	SuperAdmin := router.Group("/org").Use(middleware.JWTAuth(), middleware.SupAdmin())
	{
		// 超级管理员路由
		//uGroup.DELETE("/superadmin") //  超级管理员（组织）的注销,暂时不做支持
		SuperAdmin.POST("/credep", admin.Create) // 部门（admin）的创建
		SuperAdmin.POST("/settime", super_admin.SetTime)
		SuperAdmin.PATCH("/udporg", super_admin.Update)                   // 根据组织id更新组织信息
		SuperAdmin.GET("/getdep/detail/:depid", super_admin.GetDepDetail) // 最高权限获取部门信息,包括账号密码
		SuperAdmin.DELETE("/deldep/:depid", super_admin.Delete)           // admin，即部门的删除
		SuperAdmin.GET("/getorg", super_admin.GetOrgDep)                  // 通过id获取某一组织下的所有部门
		SuperAdmin.GET("/getdep/rough/:depid", admin.GetDetail)           // 详细获取某一个部门的信息
		SuperAdmin.GET("/getuser/info", super_admin.GetAllUserStatistic)  // 获取组织下所有报名人数及男女信息
	}
}
