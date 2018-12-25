package mssql_test

import (
	"database/sql"
	"errors"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/mssql"

	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type authenticator struct{}

func (a *authenticator) Authenticate(token string) (*root.User, error) {
	if token == "234" {
		return &root.User{No: "234", Name: "test_234"}, nil
	} else if token == "123" {
		return &root.User{No: "123", Name: "test_123"}, nil
	} else {
		return nil, errors.New("token解析错误")
	}

}

func (a *authenticator) Token(user *root.User) (string, error) {
	return "", nil
}

type MockSession struct {
	mockDB       *sql.DB
	sqlxDB       *sqlx.DB
	mock         sqlmock.Sqlmock
	mssqlSession *mssql.Session
}

func NewMockSession() (MockSession, error) {

	var mockSession MockSession

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return mockSession, err
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mssqlSession := mssql.NewSession(sqlxDB)

	authenticator := authenticator{}

	mssqlSession.SetAuthenticator(&authenticator)

	mockSession.sqlxDB = sqlxDB
	mockSession.mock = mock
	mockSession.mockDB = mockDB
	mockSession.mssqlSession = mssqlSession

	return mockSession, nil

}

func TestAuthenticate(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	// 测试token解析出现错误
	token := "1"
	mockSession.mssqlSession.SetAuthToken(token)

	u, err := mockSession.mssqlSession.Authenticate()
	if err == nil {
		t.Error("解析token时出错测试错误")
	}

	// 测试token解析正确性
	token = "123"
	// 设置token
	mockSession.mssqlSession.SetAuthToken(token)

	u, err = mockSession.mssqlSession.Authenticate()
	if err != nil {
		t.Error(err)
	}

	if u.No != "123" {
		t.Error("用户认证token即系错误")
	}

	// 测试已通过认证的session能否正确返回用户
	token = "234"
	mockSession.mssqlSession.SetAuthToken(token)
	u, err = mockSession.mssqlSession.Authenticate()
	if u.No != "123" {
		t.Errorf("用户认证失败, 期待%s, 实际%s", "123", u.No)
	}

}
