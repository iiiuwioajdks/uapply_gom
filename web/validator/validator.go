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
