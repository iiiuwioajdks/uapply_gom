package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode"
)

const (
	EMAIL = `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	PHONE = `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
)

// ValidateMobile 自定义手机验证
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 手机号码正则表达式验证
	return ValidateFunc(mobile, PHONE)
}

func ValidateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	// 邮箱正则验证
	return ValidateFunc(email, EMAIL)
}

func ValidateFunc(val string, rule string) bool {
	ok, _ := regexp.MatchString(rule, val)
	if ok {
		return true
	}
	return false
}

const (
	PwdOptNumber  uint16 = 1 << iota //数字 	 0001
	PwdOptLower                      //小写 	 0010
	PwdOptUpper                      //大写 	 0100
	PwdOptSpecial                    //特殊符号 1000
)

// VerifyPwd 密码验证
func VerifyPwd(pwd string, options uint16) bool {
	if options == 0 {
		options = PwdOptNumber | PwdOptLower | PwdOptUpper | PwdOptSpecial
	}

	if len(pwd) < 12 || len(pwd) > 18 {
		return false
	}
	// 用于记录验证结果
	var result uint16
	for _, r := range pwd {
		switch {
		case unicode.IsNumber(r):
			result = result | PwdOptNumber
		case unicode.IsLower(r):
			result = result | PwdOptLower
		case unicode.IsUpper(r):
			result = result | PwdOptUpper
		case (unicode.IsPunct(r) || unicode.IsSymbol(r)) && (r == '@' || r == '.'): //标点符号 和 字符
			result = result | PwdOptSpecial
		default:
			return false
		}
		// 比较结果和设置项
		// 当 options与result != result 表示密码字符串超出 options 范围
		if options&result != result {
			return false
		}
	}
	return true
}
