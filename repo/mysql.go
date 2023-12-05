package repo

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm/schema"
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

type Repository struct {
	*gorm.DB
}

func GetDbClient() *Repository {
	return &Repository{
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
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
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

func IsNotFound(e error) bool {

	return errors.Is(e, gorm.ErrRecordNotFound) || errors.Is(e, redis.Nil)
}

// Create 通用创建
func (rep *Repository) Create(i interface{}) error {

	return rep.DB.Create(i).Error
}

type TransactionHandle func(rep *Repository) (e error)

// Transaction 创建事务执行
func (rep *Repository) Transaction(handler TransactionHandle) error {

	tx := &Repository{DB: rep.Begin()}

	if err := handler(tx); err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
