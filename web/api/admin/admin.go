package admin

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/admin_handler"
	"uapply_go/web/middleware"
	jwt2 "uapply_go/web/models/jwt"
)

// Create 管理员（部门）的创建
func Create(c *gin.Context) {
	// 绑定参数
	var req forms.AdminReq
	c.ShouldBindJSON(&req)

	// 判断一下参数是否正确
	if req.DepartmentName == "" || req.Account == "" || req.Password == "" {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	// 获取并绑定当前的 OrganizationID, 防止篡改
	claim, ok := c.Get("claim")
	if !ok {
		zap.S().Info(claim)
	}
	claimsInfo := claim.(*jwt2.Claims)
	req.OrganizationID = claimsInfo.OrganizationID

	// 转到handler去处理
	err := admin_handler.CreateDep(&req)
	if err != nil {
		// 判断错误原因是否为重复创建部门
		if errors.Is(err, errInfo.ErrDepExist) {
			api.FailWithErr(c, api.CodeInvalidParam, err.Error())
			return
		}
		zap.S().Error("admin_handler.CreateDep()", zap.Error(err))
		// 这个可能是重复创建的索引错误，到优化阶段再改一下
		api.FailWithErr(c, api.CodeBadRequest, err.Error())
		return
	}
	// 返回给前端
	api.Success(c, "创建部门成功")
}

// Login 组织和部门都需要登录，但是组织和部门的身份由表中的role决定
// 因为登录账号密码放在department表上，所以防止admin进行
func Login(c *gin.Context) {
	var loginInfo forms.Login
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	admin, err := admin_handler.Login(context.Background(), &loginInfo)
	if errors.Is(err, errInfo.ErrUserNotFind) {
		api.Fail(c, api.CodeUserNotExist)
		return
	}
	// 记录token信息
	jwt := middleware.NewJWT()
	token, err := jwt.CreateToken(jwt2.Claims{
		Role:           int(admin.Role),
		OrganizationID: int(admin.OrganizationID),
		DepartmentID:   int(admin.DepartmentID),
	})
	if err != nil {
		api.FailWithErr(c, api.CodeSystemBusy, err.Error())
		return
	}
	// 返回token给前端
	api.Success(c, token)
}

// Update 部门更新
func Update(c *gin.Context) {
	// 绑定参数
	var req forms.AdminReq
	c.ShouldBindJSON(&req)

	claim, ok := c.Get("claim")
	if !ok {
		zap.S().Info(claim)
	}
	// 获取并绑定当前的 OrganizationID
	claimInfo := claim.(*jwt2.Claims)
	// 如果是管理员
	if claimInfo.Role == 1 && req.DepartmentID == 0 {
		api.FailWithErr(c, api.CodeInvalidParam, "组织修改时，department_id不能为空")
		return
	} else if req.DepartmentID == 0 {
		// 如果没有传给 DepartmentID ，说明是自己部门修改，而不是管理员
		req.DepartmentID = claimInfo.DepartmentID
	}
	req.OrganizationID = claimInfo.OrganizationID

	// 转到 handle 去处理
	err := admin_handler.UpdateDep(&req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("admin_handler.UpdateDep()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, "更新部门成功")
}

// GetDetail 获取某一个部门的详细信息
func GetDetail(c *gin.Context) {
	// 获取 depid
	depid, ok := c.Params.Get("depid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	//转到 handle 处理
	depInfo, err := admin_handler.GetDepDetail(depid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("admin_handler.GetDepDetail()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	// 不返回账号、密码
	api.Success(c, gin.H{
		"organization_id": depInfo.OrganizationID,
		"department_id":   depInfo.DepartmentID,
		"department_name": depInfo.DepartmentName,
	})
}

// Get 获取某一部门粗略的信息
func Get(c *gin.Context) {
	// 获取 claims
	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	//获取 depid
	depid := claimInfo.DepartmentID

	//转到 handler 处理
	depInfo, err := admin_handler.GetDepRoughDetail(depid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("admin_handler.GetDepRoughDetail()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, gin.H{
		"organization_id": depInfo.OrganizationID,
		"department_id":   depInfo.DepartmentID,
		"department_name": depInfo.DepartmentName,
	})
}

// AddExtraEnroll 直接添加一个部员
func AddExtraEnroll(c *gin.Context) {

}

// AddInterviewers 通过uid增加面试官，首先他们应该是存在且部员
func AddInterviewers(c *gin.Context) {

}

// Pass 通过某一轮面试
func Pass(c *gin.Context) {

}

// Out 在某一轮面试被淘汰
func Out(c *gin.Context) {

}

// Enroll 在某一轮面试被录取
func Enroll(c *gin.Context) {

}

// GetAllInterviewees 部门获取报名自己部门的所有用户
func GetAllInterviewees(c *gin.Context) {

}

// GetInterviewee 部门获取报名自己部门的某一个用户详细信息
func GetInterviewee(c *gin.Context) {

}

// GetUninterview 部门获取第n轮未面试成员
func GetUninterview(c *gin.Context) {

}

// Getinterviewed 部门获取第n轮已面试成员
func Getinterviewed(c *gin.Context) {

}

// GetUserEnroll 部门获取自己的通过部员
func GetUserEnroll(c *gin.Context) {

}

// GetUserInfo 获取本部门男女人数，报名人数信息
func GetUserInfo(c *gin.Context) {

}

// DeleteInterviewers 删除面试官
func DeleteInterviewers(c *gin.Context) {

}
