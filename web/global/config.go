package global

import (
	ut "github.com/go-playground/universal-translator"
	"uapply_go/web/config"
)

// Conf 所有配置都存在该变量中
var Conf config.ServerConf

var Trans ut.Translator
