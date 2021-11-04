package super_admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/handler/super_admin_handler"
)

// Create 超级管理员（组织）的创建
func Create(c *gin.Context) {
	var csa forms.CreateSAdmin
	if err := c.ShouldBindJSON(&csa); err != nil {
		zap.S().Info(err)
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	err := super_admin_handler.Create(&csa)
	if err != nil {
		zap.S().Info(err)
		api.FailWithErr(c, api.CodeInvalidParam, err.Error())
		return
	}
	api.Success(c, nil)
}

// Update 组织信息的更新
func Update(c *gin.Context) {

}

// GetOrg 获取组织信息
func GetOrg(c *gin.Context) {

}

// GetOrgDep 根据组织获取其下的附属部门
func GetOrgDep(c *gin.Context) {

}

// Delete 组织删除其下附属的某一个或多个部门
func Delete(c *gin.Context) {

}

// GetDepDetail 最高权限获取部门信息，包括账号密码，方便统一管理
func GetDepDetail(c *gin.Context) {

}
