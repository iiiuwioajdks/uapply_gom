package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"uapply_go/web/models"
)

func main() {
	dsn := fmt.Sprintf("root:root@tcp(121.41.170.94:3309)/uapply?charset=utf8&parseTime=True&loc=Local")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 不用复数
			SingularTable: true,
		},
		Logger: newLogger,
	})

	db.AutoMigrate(&models.Department{})
	db.AutoMigrate(&models.Organization{})
}
