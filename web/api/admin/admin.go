package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	claims, _ := c.Get("claim")
	claimsInfo := claims.(*jwt2.Claims)
	req.OrganizationID = claimsInfo.OrganizationID

	// 转到handler去处理
	err := admin_handler.CreateDep(&req)
	if err != nil {
		api.Fail(c, api.CodeSystemBusy)
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

	// 校验参数
	// DepartmentID 主键自增，一定大于零
	if req.DepartmentID <= 0 || req.DepartmentName == "" || req.Account == "" || req.Password == "" {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	claims, ok := c.Get("claim")
	if !ok {
		fmt.Println(claims)
	}
	// 获取并绑定当前的 OrganizationID
	claimInfo := claims.(*jwt2.Claims)
	req.OrganizationID = claimInfo.OrganizationID

	// 转到 handle 去处理
	rowsAffected, err := admin_handler.UpdateDep(&req)
	if err != nil {
		api.FailWithErr(c, api.CodeSystemBusy, err.Error())
		return
	}
	// 更新过程没有出错但 rowsAffected 等于 0 说明 DepartmentID 有误
	if rowsAffected == 0 {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	api.Success(c, "更新部门成功")
}

// GetDetail 获取某一个部门的详细信息
func GetDetail(c *gin.Context) {

}

// Get 获取某一部门粗略的信息
func Get(c *gin.Context) {

}
