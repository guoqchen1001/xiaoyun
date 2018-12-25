package jwt

import (
	"time"
	root "xiaoyun/pkg"

	jwt "github.com/dgrijalva/jwt-go"
)

// Authenticator 身份验证
type Authenticator struct {
	Secret    string // 私钥
	ExpiresAt int    // 过期时间,单位秒
	Configer  root.Configer
}

// NewAuthenticator 创建新的验证对象
func NewAuthenticator(configer root.Configer) *Authenticator {
	a := &Authenticator{}
	a.Configer = configer
	return a
}

// CustomClaims jwt实现
type CustomClaims struct {
	UserNo string `json:"user_no"`
	jwt.StandardClaims
}

// Authenticate 身份验证，实现root.Authenticator接口
func (a *Authenticator) Authenticate(tokenString string) (*root.User, error) {

	var customError root.Error
	customError.Op = "jwt.Authenticator.Authenticate"

	if err := a.init(); err != nil {
		customError.Err = err
		return nil, &customError
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		userno := claims.UserNo
		return &root.User{No: userno}, nil
	}
	return nil, err
}

// Token 获取Token
func (a *Authenticator) Token(user *root.User) (string, error) {

	var customError root.Error
	customError.Op = "jwt.Authenticator.token"

	if err := a.init(); err != nil {
		customError.Err = err
		return "", &customError
	}

	signKey := []byte(a.Secret)

	exp := time.Now().Add(time.Second * time.Duration(a.ExpiresAt))
	customClaims := CustomClaims{
		user.No,
		jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	ss, err := token.SignedString(signKey)

	if err != nil {
		return "", err
	}

	return ss, nil

}

func (a *Authenticator) init() error {

	var customError root.Error
	customError.Op = "jwt.Authenticator.Init"

	if a.Configer == nil {
		customError.Code = root.ECONFIGNOTFOUND
		return &customError
	}

	config, err := a.Configer.GetConfig()
	if err != nil {
		customError.Err = err
		return &customError
	}

	if config.Auth == nil {
		customError.Code = root.ECONFIGAUTHNOTFOUND
		return &customError
	}

	a.Secret = config.Auth.Secret
	a.ExpiresAt = config.Auth.ExpiredAt

	return nil
}
