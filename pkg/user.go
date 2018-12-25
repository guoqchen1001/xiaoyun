package root

// User 用户结构体
type User struct {
	No       string
	Name     string
	Password string
	Group    UserGroup
}

// Session 代码对服务的认证连接
type Session interface {
	SetAuthToken(token string) // 设置Token
}

// Authenticator 用户认证接口
type Authenticator interface {
	Authenticate(tokenString string) (*User, error) // 验证Token
	Token(user *User) (string, error)
}
