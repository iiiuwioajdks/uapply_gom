package errInfo

import "github.com/pkg/errors"

var (
	ErrUserNotFind = errors.New("用户不存在")
	ErrWXCode      = errors.New("微信登录code失效")
	ErrResumeExist = errors.New("简历已存在，不可重复提交")
)
