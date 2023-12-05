package main

import (
	"fmt"
	"singo/cache"
	"singo/conf"
	_ "singo/docs"
	"singo/logger"
	"singo/repo"
	"singo/server"
	"singo/util"
)

func main() {

	cache.InitRedis()
	repo.InitMysql()
	go util.InitUbq()
	// 装载路由
	r := server.NewRouter()
	if err := r.Run(fmt.Sprintf(":%d", conf.GetConfig().Server.Port)); err != nil {
		logger.Fatal("启动服务出错", err)
		return
	}
}
