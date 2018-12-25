package mysql_test

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/mysql"

	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type MockSession struct {
	mockDB       *sql.DB
	sqlxDB       *sqlx.DB
	mock         sqlmock.Sqlmock
	mysqlSession *mysql.Session
}

func NewMockSession() (MockSession, error) {

	var mockSession MockSession

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return mockSession, err
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mysqlSession := mysql.NewSession(sqlxDB)

	var authenticator Authenticator
	authenticator.authFn = func(token string) (*root.User, error) {
		return &root.User{}, nil
	}
	mysqlSession.Authenticator = &authenticator

	mockSession.sqlxDB = sqlxDB
	mockSession.mock = mock
	mockSession.mockDB = mockDB
	mockSession.mysqlSession = mysqlSession

	return mockSession, nil

}

type Authenticator struct {
	authFn  func(token string) (*root.User, error)
	tokenFn func(user *root.User) (string, error)
}

func (a *Authenticator) Authenticate(token string) (*root.User, error) {
	return a.authFn(token)
}

func (a *Authenticator) Token(user *root.User) (string, error) {
	return a.tokenFn(user)
}

func TestSession_Authenticate(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	user := root.User{
		No: "234",
	}

	var authenticator Authenticator
	authenticator.authFn = func(token string) (*root.User, error) {
		var user root.User
		if token == "123" {
			user.No = "123"
		} else if token == "234" {
			user.No = "234"
		} else {
			return nil, errors.New("auth_err")
		}

		return &user, nil
	}

	mockSession.mysqlSession.Authenticator = &authenticator

	mockSession.mysqlSession.SetAuthToken("123")

	userAuth, err := mockSession.mysqlSession.Authenticate("234")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(userAuth, &user) {
		t.Errorf("token验证出错, 期待%v, 实际%v", &user, userAuth)
	}

	userAuth, err = mockSession.mysqlSession.Authenticate("123")
	if !reflect.DeepEqual(userAuth, &user) {
		t.Errorf("token验证出错, 期待%v, 实际%v", &user, userAuth)
	}

}

func TestSession_AuthenticateError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	var authenticator Authenticator
	authenticator.authFn = func(token string) (*root.User, error) {
		return nil, errors.New("auth_err")
	}

	mockSession.mysqlSession.Authenticator = &authenticator

	_, err = mockSession.mysqlSession.Authenticate("123")
	if root.ErrorCode(err) != root.EAUTHERROR {
		t.Errorf("错误码不符合预期，预期[%s]，实际[%s]", root.EAUTHERROR, root.ErrorCode(err))
	}

}
