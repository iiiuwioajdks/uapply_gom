package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"uapply_go/web/global"
)

type CodeType int32

const (
	CodeSuccess CodeType = 1000 + iota
	CodeInvalidParam
	CodeSystemBusy
	CodeUserInfoNotExist
	CodeHasNotPower
	CodeUserNotExist
	CodeBadRequest
)

var codeMsg = map[CodeType]string{
	CodeSuccess:          "success",
	CodeInvalidParam:     "参数有误",
	CodeSystemBusy:       "系统繁忙",
	CodeUserInfoNotExist: "用户信息加载失败",
	CodeHasNotPower:      "没有权限",
	CodeUserNotExist:     "用户不存在",
	CodeBadRequest:       "错误请求",
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

func FailWithErr(c *gin.Context, code CodeType, err interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     code.Msg(),
		"errInfo": err,
	})
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		FailWithErr(c, CodeInvalidParam, err.Error())
		return
	}
	FailWithErr(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(global.Trans)))
	return
}

func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
