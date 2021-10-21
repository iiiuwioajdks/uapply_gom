package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"uapply_go/web/global"
	"uapply_go/web/initialize"
)

func main() {
	// 启动服务
	Router := initialize.Init()
	port := global.Conf.Port
	zap.S().Infof("启动服务器,端口: %d", port)
	go func() {
		Router.Run(":" + strconv.Itoa(port))
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // ctrl+c 和 kill ，对应win和linux
	<-quit
}
