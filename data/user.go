package data

import "singo/repo"

// @Description 用户序列化器
type UserReq struct {
	// 编号
	ID uint `json:"id"`
	// 用户名
	UserName string `json:"user_name"`
	// 昵称
	Nickname string `json:"nickname"`
	// 状态
	Status string `json:"status"`
	// 头像
	Avatar string `json:"avatar"`
	// 注册时间
	CreatedAt int64 `json:"created_at"`
	// 颁发Token
	Token string `json:"token,omitempty"`
	// 过期时间
	TokenExpire int64 `json:"token_expire,omitempty"`
}

// BuildUser 序列化用户
func BuildUser(user *repo.User) *UserReq {
	return &UserReq{
		ID:       user.ID,
		UserName: user.UserName,
		Nickname: user.Nickname,
		Status:   user.Status,
		Avatar:   user.Avatar,
	}
}
