# Gugo

Gugo: 用最简单的架构，实现够用的框架，服务海量用户

## 目的

本项目采用了一系列Golang中比较流行的组件，可以以本项目为基础快速搭建Restful Web API


## 组件

1. [Gin](https://github.com/gin-gonic/gin): 轻量级Web框架，自称路由速度是golang最快的
2. [GORM](https://gorm.io/index.html): ORM工具。本项目需要配合Mysql使用
3. [JWT-go](https://github.com/dgrijalva/jwt-go): Go的JWT中间件
4. [Go-Redis](https://github.com/go-redis/redis): Golang Redis客户端
5. [Gin-Cors](https://github.com/gin-contrib/cors): Gin框架提供的跨域中间件
5. [Gin-Swagger](https://github.com/swaggo/gin-swagger): Gin框架的Swagger文档
6. [Viper](https://github.com/spf13/viper): 配置读取框架
7. [zap](https://go.uber.org/zap): 非常高效的日志框架

预先实现了一些常用的代码方便参考和复用:

1. 创建了用户模型
2. 实现了```/api/v1/user/register```用户注册接口
3. 实现了```/api/v1/user/login```用户登录接口
4. 实现了```/api/v1/user/me```用户资料接口(需要登录后获取token)
5. 实现了```/api/v1/user/list```用户列表接口(需要登录后获取token)
