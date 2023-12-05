package repo

import (
	"golang.org/x/crypto/bcrypt"
	"singo/req"
)

// @Description 用户模型
type User struct {
	// 编号
	ID uint `gorm:"primarykey"`
	// 用户名
	UserName string
	// 密码加密
	PasswordDigest string
	// 昵称
	Nickname string
	// 状态
	Status string
	// 头像
	Avatar string `gorm:"size:1000"`
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 被封禁用户
	Suspend string = "suspend"
)

// GetUser 用ID获取用户
func (rep *Repository) GetUser(username string) (user *User, err error) {
	err = rep.Where("user_name = ?", username).First(&user).Error
	return
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

func (rep *Repository) GetUsers(param *req.PageUserReq) (total int64, array []*User, err error) {
	if param.UserName != "" {
		rep.DB = rep.DB.Where("user_name = ?", param.UserName)
	}

	// 查询总数
	if err = rep.Model(&User{}).Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 分页查询
	if err = rep.Offset(param.Offset()).Limit(param.Limit).Find(&array).Error; err != nil {
		return 0, nil, err
	}

	return
}
