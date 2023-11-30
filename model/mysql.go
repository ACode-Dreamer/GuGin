package model

import (
	"fmt"
	"log"
	"os"
	"singo/conf"
	"singo/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var config = conf.GetConfig()

// DbClient 数据库链接单例
var DbClient *gorm.DB

type MyDb struct {
	*gorm.DB
}

func GetDbClient() *MyDb {
	return &MyDb{
		DbClient,
	}
}

// Database 在中间件中初始化mysql链接
func InitMysql() {

	// 构建 MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)
	// 初始化GORM日志配置
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	// Error
	if dsn == "" || err != nil {
		logger.Error("mysql 连接失败: %v", err)
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("mysql 连接失败: %v", err)
		panic(err)
	}

	// 设置连接池
	// 空闲
	sqlDB.SetMaxIdleConns(10)
	// 打开
	sqlDB.SetMaxOpenConns(20)
	DbClient = db
	// 更新数据结构
	migration()
}

func migration() {
	// 自动迁移模式
	_ = DbClient.AutoMigrate(
		&User{},
	)
}
