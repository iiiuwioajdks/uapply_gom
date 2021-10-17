package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Rdb redis连接的客户端
var Rdb *redis.Client

// DB mysql连接的客户端
var DB *gorm.DB
