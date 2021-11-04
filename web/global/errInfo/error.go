package errInfo

import "github.com/pkg/errors"

var (
	ErrUserNotFind = errors.New("用户不存在")
)
