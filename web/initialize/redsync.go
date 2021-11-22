package initialize

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"uapply_go/web/global"
)

func InitRedsync() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Conf.RedisInfo.Host, global.Conf.RedisInfo.Port),
		Password: global.Conf.RedisInfo.Auth,
		DB:       global.Conf.RedisInfo.DB,
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)
	global.Rs = rs
}
