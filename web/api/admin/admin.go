package admin

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/admin_handler"
	"uapply_go/web/middleware"
	jwt2 "uapply_go/web/models/jwt"
	"uapply_go/web/models/response"
	"uapply_go/web/validator"
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
	if ok := validator.VerifyPwd(req.Password, 0); !ok {
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
	api.Success(c, gin.H{
		"token":           token,
		"department_id":   admin.DepartmentID,
		"organization_id": admin.OrganizationID,
		"role":            admin.Role,
	})
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
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	//获取 depid
	depid := claimInfo.DepartmentID

	//转到 handler 处理
	depInfo, err := admin_handler.GetDepRoughDetail(depid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeBadRequest)
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

// AddInterviewers 通过uid增加面试官，首先他应该是存在且部员
func AddInterviewers(c *gin.Context) {
	var uid forms.Interviewer
	if err := c.ShouldBindJSON(&uid); err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	err := admin_handler.AddInterviewers(claimInfo, &uid)
	if err != nil {
		if errors.Is(err, errInfo.ErrInvalidParam) {
			api.FailWithErr(c, api.CodeInvalidParam, "此人不存在你的部门或已经是你部门的面试官")
			return
		}
		if errors.Is(err, errInfo.ErrConcurrent) {
			api.Fail(c, api.CodeConcurrent)
			return
		}
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, "添加成功")
}

// Pass 通过某一轮面试
func Pass(c *gin.Context) {
	var uidsForm forms.MultiUIDForm
	num, ok := c.Params.Get("num")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	err := c.ShouldBindJSON(&uidsForm)
	if err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	//获取 depid orgid
	depid := claimInfo.DepartmentID
	orgid := claimInfo.OrganizationID
	err = admin_handler.Pass(num, orgid, depid, uidsForm)
	if err != nil {
		if errors.Is(err, errInfo.ErrInvalidParam) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("admin_handler.Pass()", zap.Error(err))
		api.FailWithErr(c, api.CodeSystemBusy, err.Error())
		return
	}

	api.Success(c, "通过面试"+num)
}

// Out 在某一轮面试被淘汰
func Out(c *gin.Context) {
	var uidsForm forms.MultiUIDForm
	// uid 数组为必填项
	if err := c.ShouldBindJSON(&uidsForm); err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 获取 claim
	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	// 获取 orgid 和 depid
	orgid := claimInfo.OrganizationID
	depid := claimInfo.DepartmentID

	err := admin_handler.Out(&uidsForm, orgid, depid)
	if err != nil {
		zap.S().Error("admin_handler.Out()", zap.Error(err))
		api.FailWithErr(c, api.CodeSystemBusy, "部分请求失败，请刷新重试")
		return
	}
	api.Success(c, "淘汰用户成功")
}

// Enroll 在某一轮面试被录取
func Enroll(c *gin.Context) {
	var uidForm forms.MultiUIDForm
	// uid 必填
	if err := c.ShouldBindJSON(&uidForm); err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	orgid := claimInfo.OrganizationID
	depid := claimInfo.DepartmentID

	err := admin_handler.Enroll(&uidForm, orgid, depid)
	if err != nil {
		zap.S().Error("admin_handler.Enroll()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, "录取用户成功")
}

// GetAllInterviewees 部门获取报名自己部门的所有用户
func GetAllInterviewees(c *gin.Context) {
	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	// 获取 orgid 和 depid
	orgid := claimInfo.OrganizationID
	depid := claimInfo.DepartmentID

	interviewees, err := admin_handler.GetAllInterviewees(depid, orgid)
	if err != nil {
		zap.S().Error("admin_handler.GetAllInterviewees()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, interviewees)
}

// GetInterviewee 部门获取报名自己部门的某一个用户详细信息
func GetInterviewee(c *gin.Context) {
	uid, ok := c.Params.Get("uid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	//获取 depid orgid
	depid := claimInfo.DepartmentID
	orgid := claimInfo.OrganizationID

	userInfo, err := admin_handler.GetInterviewee(uid, depid, orgid)
	if err != nil {
		if errors.Is(err, errInfo.ErrUserNotFind) {
			api.FailWithErr(c, api.CodeInvalidParam, "部门内未查找到此用户")
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			api.FailWithErr(c, api.CodeInvalidParam, "用户表中无此用户")
			return
		}
		zap.S().Error("admin.GetInterviewee", err)
		api.Fail(c, api.CodeSystemBusy)
		return
	}

	api.Success(c, userInfo)

}

// GetUninterview 部门获取第n轮未面试成员
func GetUninterview(c *gin.Context) {

}

// GetInterviewed 部门获取第n轮已面试成员
func GetInterviewed(c *gin.Context) {
	num, ok := c.Params.Get("num")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	//获取 depid orgid
	depid := claimInfo.DepartmentID
	orgid := claimInfo.OrganizationID

	// handler
	interviewers, err := admin_handler.GetInterviewed(num, orgid, depid)
	if err != nil {
		if errors.Is(err, errInfo.ErrInvalidParam) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("admin_handler.GetInterviewed()", zap.Error(err))
		api.FailWithErr(c, api.CodeSystemBusy, err.Error())
		return
	}

	var resp []*response.Interviewed

	for _, item := range interviewers {
		resp = append(resp, &response.Interviewed{
			UID:  item.UID,
			Name: item.Name,
		})
	}

	api.Success(c, resp)

}

// GetUserEnroll 部门获取自己的通过部员
func GetUserEnroll(c *gin.Context) {

}

// GetUserInfo 获取本部门男女人数，报名人数信息
func GetUserInfo(c *gin.Context) {
	// 获取claims
	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt2.Claims)
	// 获取depid 和 orgid
	depid := claimInfo.DepartmentID
	orgid := claimInfo.OrganizationID
	rsp, err := admin_handler.GetUserInfo(depid, orgid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.FailWithErr(c, api.CodeUserNotExist, "本部门暂无人报名")
			return
		}
		api.Fail(c, api.CodeSystemBusy)
		return
	}

	api.Success(c, rsp)
}

// DeleteInterviewers 删除面试官
func DeleteInterviewers(c *gin.Context) {

}

// SetTime 设置报名开始和结束时间,部门可以通过这个接口设置报名时间
func SetTime(c *gin.Context) {
	var t forms.Time
	err := c.ShouldBindJSON(&t)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	if time.Now().Unix() > t.End || t.End < t.Start {
		api.FailWithErr(c, api.CodeBadRequest, "无效的设置")
		return
	}

	claim, ok := c.Get("claim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimsInfo := claim.(*jwt2.Claims)

	err = admin_handler.SetTime(claimsInfo.DepartmentID, &t)
	if err != nil {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	api.Success(c, nil)
}

func GetInterviewers(c *gin.Context) {

}

func SMS(c *gin.Context) {
	var sms forms.MultiUIDForm
	if err := c.ShouldBindJSON(&sms); err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	ok, depid, orgid := GetID(c)
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	phoneInfo, err := admin_handler.GetPhones(depid, orgid, &sms)
	if err != nil {
		if errors.Is(err, errInfo.ErrInvalidParam) {
			api.FailWithErr(c, api.CodeInvalidParam, "有用户没有对应的手机号或者有用户不为你部门的成员，请检查选中人员，若多次出现请联系相关人员")
			return
		}
		zap.S().Error("api.admin.SMS.Func(GetPhones):", err)
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	// 根据手机号码和uid发送短信和面试密钥
	err = admin_handler.SendSMS(sms.Type, phoneInfo)
	if err != nil {
		if errors.Is(err, errInfo.ErrInvalidParam) {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
		zap.S().Error("api.admin.SMS.Func(SendSMS):", err)
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, phoneInfo)
}

func GetID(c *gin.Context) (bool, int, int) {
	// 获取claims
	claim, ok := c.Get("claim")
	if !ok {
		return false, 0, 0
	}
	claimInfo := claim.(*jwt2.Claims)
	// 获取depid 和 orgid
	depid := claimInfo.DepartmentID
	orgid := claimInfo.OrganizationID
	return true, depid, orgid
}
