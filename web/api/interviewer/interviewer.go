package interviewer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"uapply_go/web/api"
	"uapply_go/web/global/errInfo"
	"uapply_go/web/handler/inter_handler"
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

}

func GetDepInfo(c *gin.Context) {

}

func Evaluate(c *gin.Context) {

}
