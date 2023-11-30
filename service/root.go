package service

import (
	"singo/cache"
	"singo/model"
)

func rep() *model.MyDb {
	return model.GetDbClient()
}

func redis() *cache.MyRedis {
	return cache.GetRedisClient()
}
