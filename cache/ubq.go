package cache

import (
	"singo/logger"
	"strconv"
	"time"
)

func (rep *MyRedis) SetUbqToken(token string, expire time.Duration) (err error) {
	err = rep.Set("ubq:token", token, expire).Err()
	return
}

func (rep *MyRedis) GetUbqToken() string {
	return rep.Get("ubq:token").Val()
}

func (rep *MyRedis) GetUbqHealth() int {
	energy, err := rep.Get("ubq:health").Int()
	if err != nil {
		logger.Error("获取体力能量报错", err)
		return 199
	}
	return energy
}

func (rep *MyRedis) GetRepeatBalance() (balance int, err error) {

	cmd := rep.Get("game:repeat")

	if err = cmd.Err(); err != nil {

		return
	}

	balance, err = strconv.Atoi(cmd.Val())

	return
}
