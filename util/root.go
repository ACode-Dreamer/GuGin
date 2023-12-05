package util

import (
	"singo/cache"
	"singo/repo"
)

func rep() *repo.Repository {
	return repo.GetDbClient()
}

func redis() *cache.MyRedis {
	return cache.GetRedisClient()
}
