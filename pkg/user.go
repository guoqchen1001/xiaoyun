package root

// User 用户结构体
type User struct {
	No       string `json:"no"`
	Name     string `json:"name"`
	Password string `json:"password"`
	GroupNo  string `json:"group_no"`
}

// Credentials 登录证书
type Credentials struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
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

// UserService 用户服务
type UserService interface {
	Login(c Credentials) (string, error)
	CreateUser(user *User) error
}
