package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CodeType int32

const (
	CodeSuccess CodeType = 1000 + iota
	CodeInvalidParam
	CodeSystemBusy
)

var codeMsg = map[CodeType]string{
	CodeInvalidParam: "参数有误",
	CodeSystemBusy:   "系统繁忙",
}

func (c CodeType) Msg() string {
	if v, ok := codeMsg[c]; ok {
		return v
	}
	return "状态码获取出错"
}
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.Msg(),
		"data": data,
	})
}

func Fail(c *gin.Context, code CodeType) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.Msg(),
	})
}
