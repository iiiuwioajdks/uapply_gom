package main

import (
	"fmt"
	"go.uber.org/zap"
	"uapply_go/user_web/global"
	"uapply_go/user_web/initialize"
)

func main() {
	// 配置初始化
	err := initialize.ViperInit()
	if err != nil {
		panic(err)
	}
	// 日志初始化
	err = initialize.InitLogger(global.Conf.LoggerInfo.Mode)
	if err != nil {
		panic(err)
	}
	// 路由初始化
	Router := initialize.InitRouter()
	// mysql初始化
	err = initialize.InitMysql()
	if err != nil {
		panic(err)
	}
	// redis初始化
	err = initialize.InitRedis()
	if err != nil {
		panic(err)
	}
	// 启动服务
	port := global.Conf.Port
	zap.S().Infof("启动服务器,端口: %d", port)
	err = Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
