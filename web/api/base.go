package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func InvalidParam(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": "参数有误",
		"err": err.Error(),
	})
}

func HandlerErr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"msg": "系统繁忙",
	})
}
