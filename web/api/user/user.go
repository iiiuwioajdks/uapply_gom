package user

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/user_handler"
	jwt2 "uapply_go/web/models/jwt"
)

// Login 微信小程序用户端登录
func Login(c *gin.Context) {
	code := c.Param("code")
	token, err := user_handler.Login(code)
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
	api.Success(c, token)
}

func SaveResume(c *gin.Context) {
	// 绑定参数
	var req forms.UserInfoReq
	c.ShouldBindJSON(&req)

	// 未加 binding 要求，需要校验参数
	// 楼号(Address)和专业(Major)是可选项，不用校验
	if req.Name == "" || req.StuNum == "" || req.Phone == "" || req.Email == "" || req.Intro == "" || (req.Sex != 0 && req.Sex != 1) {
		api.Fail(c, api.CodeInvalidParam)
		return
	}
	// 获取 wxClaim
	wxClaim, ok := c.Get("wxClaim")
	if !ok {
		zap.S().Info(wxClaim)
	}
	wxClaimInfo := wxClaim.(*jwt2.WXClaims)
	// 获取 OpenID
	openID := wxClaimInfo.Openid
	// 获取 UID
	uid, err := user_handler.GetUID(openID)
	if err != nil {
		// 没有获取到 UID
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeUserNotExist)
			return
		}
		zap.S().Error("user_handler.GetUID()", zap.Error(err))
		api.Fail(c, api.CodeSystemBusy)
		return
	}
	req.UID = uid

	//转到 handler 处理
	err = user_handler.SaveResume(&req)
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

}

func ClearText(c *gin.Context) {

}
