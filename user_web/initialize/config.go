package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"uapply_go/user_web/global"
)

// ViperInit viper配置文件初始化
func ViperInit() error {
	Conf := &global.Conf
	viper.SetConfigFile("user_web/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "setting init error")
	}
	if err := viper.Unmarshal(Conf); err != nil {
		return errors.Wrap(err, "unmarshal init error")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return nil
}
