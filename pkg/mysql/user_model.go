package mysql

import root "xiaoyun/pkg"

// UserNodel 用户模型
type UserNodel struct {
	No       string `db:"no"`
	Password string `db:"password"`
	Name     string `db:"name"`
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

	return model

}
