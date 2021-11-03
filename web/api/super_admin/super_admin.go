package super_admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"uapply_go/web/api"
	"uapply_go/web/forms"
	"uapply_go/web/handler/super_admin_handler"
)

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
