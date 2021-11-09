package router

import (
	"github.com/gin-gonic/gin"
	"uapply_go/web/api/admin"
	"uapply_go/web/api/super_admin"
)

func InitCommonRouter(router *gin.RouterGroup) {
	cGroup := router.Group("/common")
	{
		cGroup.POST("/superadmin", super_admin.Create)      // 超级管理员（组织）的创建,该接口不对外暴露，由后端人员操作
		cGroup.GET("/org/:orgid", super_admin.GetOrg)       // 通过id获取某一组织的信息,这个不需要权限，想获取就获取
		cGroup.GET("/alldep/:orgid", super_admin.GetOrgDep) // 通过id获取某一组织下的所有部门
		cGroup.GET("/onedep/:depid", admin.GetDetail)       // 获取某一个部门的详细信息
		cGroup.POST("/login", admin.Login)                  // 部门或组织登录
	}
}
