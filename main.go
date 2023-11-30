package main

import (
	"fmt"
	"singo/cache"
	"singo/conf"
	_ "singo/docs"
	"singo/logger"
	"singo/model"
	"singo/server"
)

func main() {

	cache.InitRedis()
	model.InitMysql()
	// 装载路由
	r := server.NewRouter()
	if err := r.Run(fmt.Sprintf(":%d", conf.GetConfig().Server.Port)); err != nil {
		logger.Fatal("启动服务出错", err)
		return
	}
}
