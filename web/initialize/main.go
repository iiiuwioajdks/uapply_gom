package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"uapply_go/web/global"
)

func Init() *gin.Engine {
	// 配置初始化
	err := ViperInit()
	if err != nil {
		panic(err)
	}
	fmt.Println(global.Conf)
	// 日志初始化
	err = InitLogger(global.Conf.LoggerInfo.Mode)
	if err != nil {
		panic(err)
	}
	// 路由初始化
	Router := InitRouter()
	// 注册验证器
	InitValidators()
	InitTrans("zh")
	// mysql初始化
	err = InitMysql()
	if err != nil {
		panic(err)
	}
	// redis初始化
	err = InitRedis()
	if err != nil {
		panic(err)
	}
	return Router
}
