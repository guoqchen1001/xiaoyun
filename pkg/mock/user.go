package mock

import (
	root "xiaoyun/pkg"
)

// UserService 用户服务
type UserService struct {
	LoginFn      func(credentials root.Credentials) (string, error)
	CreateUserFn func(user *root.User) error
	LoginInvoked bool
}

// Login 登录
func (u *UserService) Login(credentials root.Credentials) (string, error) {
	u.LoginInvoked = true
	return u.LoginFn(credentials)
}

// CreateUser 创建用户
func (u *UserService) CreateUser(user *root.User) error {
	u.LoginInvoked = true
	return u.CreateUserFn(user)
}
