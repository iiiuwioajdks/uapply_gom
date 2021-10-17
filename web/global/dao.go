package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Rdb redis连接的客户端
var Rdb *redis.Client

// DB mysql连接的客户端
var DB *gorm.DB

// GetDB 使用这个使得各个db之间的操作互不影响
func GetDB(ctx context.Context) *gorm.DB {
	return DB.WithContext(ctx)
}
