package errInfo

import "github.com/pkg/errors"

var (
	ErrUserNotFind    = errors.New("用户不存在")
	ErrWXCode         = errors.New("微信登录code失效")
	ErrResumeExist    = errors.New("简历已存在，不可重复提交")
	ErrResumeNotExist = errors.New("简历不存在，请先填写简历")
	ErrInvalidParam   = errors.New("参数错误")
	ErrReRegister     = errors.New("一个组织只能报名一个部门")
)
