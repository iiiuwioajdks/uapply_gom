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

// GetNewDB 使用一个新Session创建一个子db，使得各个db之间的操作互不影响
func GetNewDB(ctx context.Context) *gorm.DB {
	return DB.WithContext(ctx)
}

// https://gorm.io/zh_CN/docs/method_chaining.html 进行详细讲解
