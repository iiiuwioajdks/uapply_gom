package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"uapply_go/web/forms"
	"uapply_go/web/handler"
)

func CSuperAdmin(c *gin.Context) {
	var csa forms.CreateSAdmin
	if err := c.ShouldBindJSON(&csa); err != nil {
		zap.S().Info(err)
		InvalidParam(c, err)
		return
	}
	err := handler.CreateSAdmin(&csa)
	if err != nil {
		zap.S().Info(err)
		HandlerErr(c)
		return
	}
	Success(c, nil)
}
