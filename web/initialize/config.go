package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"uapply_go/web/global"
)

// ViperInit viper配置文件初始化
func ViperInit() error {
	workdir, _ := os.Getwd()
	Conf := &global.Nacos
	viper.SetConfigFile(workdir + "/config.yaml")
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

	// 连接nacos
	sc := []constant.ServerConfig{
		{
			IpAddr: Conf.Host,
			Port:   uint64(Conf.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         Conf.NamespaceId, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return err
	}

	// 获取文件内容
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: Conf.DataId,
		Group:  Conf.Group,
	})
	if err != nil {
		zap.S().Error(err)
		return err
	}
	_ = json.Unmarshal([]byte(content), &global.Conf)
	zap.S().Info(global.Conf)
	// 监听变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: Conf.DataId,
		Group:  Conf.Group,
		OnChange: func(namespace, group, dataId, data string) {
			content, _ := client.GetConfig(vo.ConfigParam{
				DataId: Conf.DataId,
				Group:  Conf.Group,
			})
			_ = json.Unmarshal([]byte(content), &global.Conf)
			Init()
			zap.S().Info(global.Conf)
		},
	})
	return nil
}
