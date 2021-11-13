package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/user_handler"
	jwt2 "uapply_go/web/models/jwt"
	validator2 "uapply_go/web/validator"
)

// Login 微信小程序用户端登录
func Login(c *gin.Context) {
	code := c.Param("code")
	token, uid, err := user_handler.Login(code)
	if err != nil {
		zap.L().Error("wxapp1 login error", zap.Error(err))
		log.Println(err)
		if errors.Is(err, errInfo.ErrWXCode) {
			api.FailWithErr(c, api.CodeInvalidParam, errInfo.ErrWXCode.Error())
			return
		}
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, gin.H{
		"token": token,
		"uid":   uid,
	})
}

func SaveResume(c *gin.Context) {
	// 绑定参数
	var req forms.UserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErr(c, api.CodeInvalidParam, err.Error())
		return
	}
	// 获取 wxClaim
	wxClaim, ok := c.Get("wxClaim")
	if !ok {
		zap.S().Info(wxClaim)
	}
	wxClaimInfo := wxClaim.(*jwt2.WXClaims)
	// 获取 UID
	req.UID = wxClaimInfo.UID

	//转到 handler 处理
	err := user_handler.SaveResume(&req)
	if err != nil {
		zap.S().Error("user_handler.SaveResume()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, "保存简历成功")
}

func SaveTmpResume(c *gin.Context) {

}

func Register(c *gin.Context) {
	// 绑定参数
	var regInfo forms.UserRegisterInfo
	if err := c.ShouldBindJSON(&regInfo); err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 获取wxClaim
	wxClaim, ok := c.Get("wxClaim")
	if !ok {
		zap.S().Info(wxClaim)
	}
	wxClaimInfo := wxClaim.(*jwt2.WXClaims)

	// 绑定UID
	regInfo.UID = wxClaimInfo.UID

	// 转移handler
	err := user_handler.Register(&regInfo)
	if err != nil {
		zap.S().Error("user_handler.Register()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}

	api.Success(c, "用户报名成功")

}

func GetRegInfo(c *gin.Context) {

}

func GetRegStatus(c *gin.Context) {

}

func GetTmpResume(c *gin.Context) {

}

func GetResume(c *gin.Context) {

}

func UpdateResume(c *gin.Context) {
	var req forms.UserResumeInfo
	//绑定参数
	err := c.ShouldBindJSON(&req)
	if err != nil {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	// 需要更新电话时校验电话
	if req.Phone != "" {
		validator := validator.New()
		validator.RegisterValidation("mobile", validator2.ValidateMobile)
		err := validator.Var(req.Phone, "mobile")
		if err != nil {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
	}
	// 需要更新邮箱时校验邮箱
	if req.Email != "" {
		validator := validator.New()
		validator.RegisterValidation("email", validator2.ValidateEmail)
		err := validator.Var(req.Email, "email")
		if err != nil {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
	}
	// 需要更新性别时校验性别
	if req.Sex != 0 {
		if req.Sex != 1 && req.Sex != 2 {
			api.Fail(c, api.CodeInvalidParam)
			return
		}
	}

	//获取 claim
	claim, ok := c.Get("wxClaim")
	if !ok {
		zap.S().Info(claim)
	}
	wxClaim := claim.(*jwt2.WXClaims)
	// 绑定 uid
	req.UID = wxClaim.UID
	// 转到 handler 处理
	err = user_handler.UpdateResume(&req)
	if err != nil {
		zap.S().Error("user_handler.UpdateResume()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	api.Success(c, "更新简历成功")
}

func ClearText(c *gin.Context) {

}
