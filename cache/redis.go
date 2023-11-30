package cache

import (
	"singo/conf"
	"singo/logger"
	"strconv"

	"github.com/go-redis/redis"
)

var config = conf.GetConfig()

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

type MyRedis struct {
	*redis.Client
}

func GetRedisClient() *MyRedis {
	return &MyRedis{
		RedisClient,
	}
}

// Redis 在中间件中初始化redis链接
func InitRedis() {
	db, _ := strconv.ParseUint(config.Redis.Db, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:       config.Redis.Address,
		Password:   config.Redis.Password,
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping().Result()

	if err != nil {
		logger.Panic("连接Redis不成功", err)
	}

	RedisClient = client
}
