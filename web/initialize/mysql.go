package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"uapply_go/web/global"
)

var err error

// InitMysql gorm 连接mysql
func InitMysql() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		global.Conf.MysqlInfo.UserName,
		global.Conf.MysqlInfo.Password,
		global.Conf.MysqlInfo.Host,
		global.Conf.MysqlInfo.Port,
		global.Conf.MysqlInfo.DBName,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 不用复数
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err == nil {
		zap.S().Info("mysql 初始化成功")
	}
	return err
}
