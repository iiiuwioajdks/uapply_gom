package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redsync/redsync/v4"
	"uapply_go/web/config"
)

// Conf 所有配置都存在该变量中
var Conf config.ServerConf

var Trans ut.Translator

var Rs *redsync.Redsync
