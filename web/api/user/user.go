package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"uapply_go/web/api"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/user_handler"
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
