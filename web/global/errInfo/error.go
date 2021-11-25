package errInfo

import "github.com/pkg/errors"

var (
	ErrUserNotFind    = errors.New("用户不存在")
	ErrWXCode         = errors.New("微信登录code失效")
	ErrResumeExist    = errors.New("简历已存在，不可重复提交")
	ErrResumeNotExist = errors.New("简历不存在，请先填写简历")
	ErrInvalidParam   = errors.New("参数错误")
	ErrReRegister     = errors.New("一个组织只能报名一个部门")
	ErrDepExist       = errors.New("部门已存在，不能重复创建")
	ErrUserMatch      = errors.New("用户报名组织或部门和面试不匹配")
	ErrCNotReg        = errors.New("该时间段无法报名该部门")
	ErrSystem         = errors.New("系统繁忙")
	ErrConcurrent     = errors.New("并发错误")
	ErrInvalidUIDS    = errors.New("存在不正确的UID")
)
