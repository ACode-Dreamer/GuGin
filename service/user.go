package service

import (
	"github.com/dgrijalva/jwt-go"
	"singo/conf"
	"singo/data"
	"singo/logger"
	"singo/middleware"
	"singo/repo"
	"singo/req"
	"time"
)

// @Description 用户注册请求
type UserRegisterReq struct {
	// 昵称
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	// 用户名
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	// 密码
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
	// 密码
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

// valid 验证表单
func (service *UserRegisterReq) valid() *data.Response {
	if service.PasswordConfirm != service.Password {
		return data.NewErrorResponse(40001, "两次输入的密码不相同")
	}

	count := int64(0)
	rep().Model(&repo.User{}).Where("nickname = ?", service.Nickname).Count(&count)
	if count > 0 {
		return data.NewErrorResponse(40001, "昵称被占用")
	}

	count = 0
	rep().Model(&repo.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count > 0 {
		return data.NewErrorResponse(40001, "用户名已经注册")
	}

	return nil
}

// Register 用户注册
func Register(service *UserRegisterReq) *data.Response {
	user := repo.User{
		Nickname: service.Nickname,
		UserName: service.UserName,
		Status:   repo.Active,
	}

	// 表单验证
	if resp := service.valid(); resp != nil {
		return resp
	}

	// 加密密码
	if err := user.SetPassword(service.Password); err != nil {
		return data.NewErrorResponse(
			data.CodeEncryptError,
			"密码加密失败",
		)
	}

	// 创建用户
	if err := rep().Create(&user); err != nil {
		logger.Error("创建用户错误", err)
		return data.NewErrorResponse(20001, "注册失败")
	}

	return data.NewDataResponse(data.BuildUser(&user))
}

// @Description 管理用户登录的请求
type UserLoginReq struct {
	// 用户名
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	// 密码
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
}

// Login 用户登录函数
func Login(service *UserLoginReq) *data.Response {
	user, err := rep().GetUser(service.UserName)
	if err != nil {
		logger.Error("查询用户错误", err)
		return data.NewErrorResponse(20002, "查询用户失败")
	}

	if !user.CheckPassword(service.Password) {
		return data.NewErrorResponse(20003, "账号密码错误")
	}

	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &middleware.Claims{
		Username: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(conf.GetConfig().Server.Secret))
	if err != nil {
		logger.Error("颁发Token错误", err)
		return data.NewErrorResponse(10000, "颁发Token错误")
	}

	resp := data.BuildUser(user)
	resp.Token = tokenString
	resp.TokenExpire = 2 * time.Hour.Milliseconds()

	if err = redis().SetToken(user.UserName, resp.Token); err != nil {
		logger.Error("存储Token错误", err)
		return data.NewErrorResponse(10000, "颁发Token错误")
	}
	return data.NewDataResponse(resp)
}

func Me(username string) *data.Response {
	user, err := rep().GetUser(username)
	if err != nil {
		return data.NewErrorResponse(20004, "查询个人信息错误")
	}
	return data.NewDataResponse(user)
}

func GetAllUsers(param *req.PageUserReq) *data.Response {
	total, array, err := rep().GetUsers(param)
	if err != nil {
		return nil
	}
	return data.NewPageResponse(total, array)
}
