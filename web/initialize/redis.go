package initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
	"uapply_go/web/global"
)

// InitRedis 连接redis
func InitRedis() error {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Conf.RedisInfo.Host, global.Conf.RedisInfo.Port),
		Password: global.Conf.RedisInfo.Auth,
		DB:       global.Conf.RedisInfo.DB,
	})

	ctx, channel := context.WithTimeout(context.Background(), 5*time.Second)
	defer channel()

	_, err = global.Rdb.Ping(ctx).Result()
	if err == nil {
		zap.S().Info("redis 初始化成功")
	}
	return err
}
