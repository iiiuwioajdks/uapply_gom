package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"uapply_go/web/global"
	"uapply_go/web/initialize"
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

	go func() {
		Router.Run(":9090")
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // ctrl+c 和 kill ，对应win和linux
	<-quit
}
