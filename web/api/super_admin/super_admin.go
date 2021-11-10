package super_admin

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/handler/super_admin_handler"
	"uapply_go/web/models/jwt"
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
	//获取部门id
	depid, ok := c.Params.Get("depid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 组织删除部门，需要先判断一下是否属于本组织，防止篡改
	claim, ok := c.Get("claim")
	if !ok {
		fmt.Println(claim)
	}
	claimInfo := claim.(*jwt.Claims)
	//转移到handler处理
	err := super_admin_handler.DeleteDepartment(depid, claimInfo.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("super_admin_handler.DeleteDepartment()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, depid)
}

// GetDepDetail 最高权限获取部门信息，包括账号密码，方便统一管理
func GetDepDetail(c *gin.Context) {

}
