package interviewer

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"uapply_go/web/api"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/inter_handler"
	"uapply_go/web/models/jwt"
)

// Login 面试官小程序登录
func Login(c *gin.Context) {
	code := c.Param("code")
	token, uid, err := inter_handler.Login(code)
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
		"uid":   uid,
		"token": token,
	})
}

func Position(c *gin.Context) {

}

func GetUser(c *gin.Context) {
	// 获取参数
	userUid, ok := c.Params.Get("uid")
	if !ok {
		api.Fail(c, api.CodeInvalidParam)
		return
	}

	// 获取token
	claim, ok := c.Get("wxClaim")
	if !ok {
		api.Fail(c, api.CodeBadRequest)
		return
	}
	claimInfo := claim.(*jwt.WXClaims)
	interUid := claimInfo.UID

	// handle处理
	useMsg, err := inter_handler.GetUser(userUid, interUid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Fail(c, api.CodeUserNotExist)
			return
		}
		if errors.Is(err, errInfo.ErrUserMatch) {
			api.Fail(c, api.CodeHasNotPower)
			return
		}
		api.Fail(c, api.CodeSystemBusy)
		zap.S().Error("inter_handler.GetUser(userUid,interUid)", zap.Error(err))
		return
	}
	api.Success(c, useMsg)
}

func GetDepInfo(c *gin.Context) {

}

func Evaluate(c *gin.Context) {

}
