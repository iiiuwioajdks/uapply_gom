package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 自定义手机验证
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 手机号码正则表达式验证
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if ok {
		return true
	}
	return false
}

func ValidateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	// 邮箱正则验证
	ok, _ := regexp.MatchString(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`, email)
	if ok {
		return true
	}
	return false
}
