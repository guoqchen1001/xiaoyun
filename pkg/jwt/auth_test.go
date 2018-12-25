package jwt_test

import (
	"testing"
	"time"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/jwt"
)

func TestJwt_Token(t *testing.T) {

	user := root.User{
		No:   "0122",
		Name: "测试",
	}
	s := jwt.Authenticator{}

	s.Secret = "kmtech"

	tokenString, err := s.Token(&user)
	if err != nil {
		t.Error(err)
	}

	t.Log(tokenString)

}

func TestJwt_Authenticate(t *testing.T) {

	user := &root.User{
		No:   "0122",
		Name: "测试",
	}
	a := jwt.Authenticator{}

	a.Secret = "kmtech"

	tokenString, err := a.Token(user)
	if err != nil {
		t.Error(err)
	}

	user, err = a.Authenticate(tokenString)
	if err != nil {
		t.Error(err)
		return
	}

	if user.No != "0122" {
		t.Error("解析不正确")
	}

	t.Log(user)
}

func TestJwt_Authenticate_wrong(t *testing.T) {

	tokenStr := "123456789"

	a := jwt.Authenticator{}
	a.Secret = "kmtech"

	_, err := a.Authenticate(tokenStr)
	if err == nil {
		t.Error(err)
	}

}

func TestJwt_Authenticate_Expired(t *testing.T) {

	user := &root.User{
		No:   "0122",
		Name: "测试",
	}
	a := jwt.Authenticator{}

	a.Secret = "kmtech"
	a.ExpiresAt = 2

	tokenString, err := a.Token(user)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Duration(a.ExpiresAt*2) * time.Second)

	user, err = a.Authenticate(tokenString)
	if err != nil {
		t.Error(err)
		return
	}

}
