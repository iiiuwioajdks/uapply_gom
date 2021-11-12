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
	// 绑定参数
	var req forms.UpdateSAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Info(err)
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	claim, ok := c.Get("claim")
	if !ok {
		zap.S().Info(claim)
	}
	// 获取并绑定当前的 OrganizationID
	claimInfo := claim.(*jwt.Claims)
	req.DepartmentID = claimInfo.DepartmentID
	req.OrganizationID = claimInfo.OrganizationID
	err := super_admin_handler.Update(&req)
	if err != nil {
		zap.S().Info(err)
		api.FailWithErr(c, api.CodeInvalidParam, err.Error())
		return
	}
	api.Success(c, "组织信息更新成功")

}

// GetOrg 获取组织信息
func GetOrg(c *gin.Context) {
	// 获取组织id
	orgid, ok := c.Params.Get("orgid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 转移handler处理
	orgInfo, err := super_admin_handler.GetOrganizationInfo(orgid)
	if err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	api.Success(c, orgInfo)
}

// GetOrgDep 根据组织获取其下的附属部门
func GetOrgDep(c *gin.Context) {
	// 获取组织id
	orgid, ok := c.Params.Get("orgid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 转移handler处理
	deps, err := super_admin_handler.GetOrgDepartments(orgid)
	if err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	api.Success(c, deps)

}

// Delete 组织删除其下附属的某一个或多个部门
func Delete(c *gin.Context) {
	//获取部门id
	depid, ok := c.Params.Get("depid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 前面的中间件只是会判断他已登录和是超级管理者
	// 组织删除部门，需要先判断一下是否属于本组织，防止篡改
	claim, ok := c.Get("claim")
	if !ok {
		fmt.Println(claim)
	}
	claimInfo := claim.(*jwt.Claims)
	//转移到handler处理
	err := super_admin_handler.DeleteDepartment(depid, claimInfo.OrganizationID)
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
	// 获取部门的ID
	depid, ok := c.Params.Get("depid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 获取token信息
	claim, ok := c.Get("claim")
	if !ok {
		// 调试信息 可删
		fmt.Println(claim)
	}
	claimInfo := claim.(*jwt.Claims)

	// 转移到handler处理
	info, err := super_admin_handler.ShowConcreteDepartInfo(depid, claimInfo.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("super_admin_handler.ShowConcreteDepartInfo()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	// 返回
	api.Success(c, info)
}
