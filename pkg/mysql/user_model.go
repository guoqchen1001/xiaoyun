package mysql

import (
	"errors"
	root "xiaoyun/pkg"
)

// UserNodel 用户模型
type UserNodel struct {
	No       string `db:"no"`
	Password string `db:"password"`
	Name     string `db:"name"`
	GroupNo  string `db:"group_no"`
}

func (u UserNodel) toUser() root.User {
	var user root.User
	user.No = u.No
	user.Name = u.Name
	user.Password = u.Password
	return user
}

func toUserModel(user *root.User) UserNodel {

	var model UserNodel

	model.No = user.No
	model.Name = user.Name
	model.Password = user.No + user.Password
	model.GroupNo = user.GroupNo

	return model

}

// Validate 验证用户模型有效性
func (u *UserNodel) Validate() error {

	if u.No == "" {
		return errors.New("缺少用户编码字段")
	}

	if u.Name == "" {
		return errors.New("缺少用户名称字段")
	}

	if u.GroupNo == "" {
		return errors.New("缺少用户组字段")
	}

	return nil

}
