package root

// UserGroup 组
type UserGroup struct {
	No   string
	Name string
}

// UserGroupService 组服务
type UserGroupService interface {
	GetUserGroup(no string) (*UserGroup, error)
	GetUserGroups() ([]*UserGroup, error)
	CreateUserGroup(group *UserGroup) error
	UpdateUserGroup(group *UserGroup) error
	DeleteUserGroup(no string) error
}
