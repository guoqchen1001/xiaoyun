package mysql_test

import (
	"database/sql"
	"errors"
	"testing"
	root "xiaoyun/pkg"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestUser_CreateOK(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	session.mock.ExpectQuery(`^SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnError(sql.ErrNoRows)
	session.mock.ExpectBegin()
	stmt := session.mock.ExpectPrepare("INSERT INTO users")

	result := sqlmock.NewResult(0, 1)
	stmt.ExpectExec().WillReturnResult(result)
	session.mock.ExpectCommit()

	service := session.mysqlSession.UserService()
	err = service.CreateUser(&user)
	if err != nil {
		t.Error(err)
	}

}

func TestUser_CreateQueryError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	err = errors.New("查询错误")
	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnError(err)

	service := session.mysqlSession.UserService()

	err = service.CreateUser(&user)

	if root.ErrorCode(err) != root.EDBQUERYERROR {
		t.Errorf("错误值返回不符合预期, 预期[%s], 实际[%s]", root.EDBQUERYERROR, root.ErrorCode(err))
	}

}

func TestUser_CreateConflict(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	columns := []string{
		"no", "name", "password",
	}

	rows := sqlmock.NewRows(columns).
		AddRow("0001", "测试", "0001")

	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnRows(rows)

	service := session.mysqlSession.UserService()

	err = service.CreateUser(&user)
	if root.ErrorCode(err) != root.ECONFLICT {
		t.Errorf("错误值返回不符合预期, 预期[%s], 实际[%s], 错误信息: %s", root.ECONFLICT, root.ErrorCode(err), err.Error())
	}

}

func TestUser_CreateBeginError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnError(sql.ErrNoRows)

	err = errors.New("BeginTrans Error")
	session.mock.ExpectBegin().WillReturnError(err)

	service := session.mysqlSession.UserService()
	err = service.CreateUser(&user)
	if root.ErrorCode(err) != root.EDBBEGINERROR {
		t.Errorf("错误值返回不符合预期, 预期[%s], 实际[%s], 错误信息: %s", root.EDBBEGINERROR, root.ErrorCode(err), err.Error())
	}

}

func TestUser_CreatePrepareError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnError(sql.ErrNoRows)

	err = errors.New("Prepare Error")
	session.mock.ExpectBegin()
	session.mock.ExpectPrepare("INSERT INTO users").WillReturnError(err)

	service := session.mysqlSession.UserService()
	err = service.CreateUser(&user)
	if root.ErrorCode(err) != root.EDBPREPAREERROR {
		t.Errorf("错误值返回不符合预期, 预期[%s], 实际[%s], 错误信息: %s", root.EDBPREPAREERROR, root.ErrorCode(err), err.Error())
	}

}

func TestUser_CreateExecError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnError(sql.ErrNoRows)
	err = errors.New("Exec Error")
	session.mock.ExpectBegin()
	stmt := session.mock.ExpectPrepare("INSERT INTO users")
	stmt.ExpectExec().WillReturnError(err)
	session.mock.ExpectRollback()

	service := session.mysqlSession.UserService()
	err = service.CreateUser(&user)
	if root.ErrorCode(err) != root.EDBEXECERROR {
		t.Errorf("错误值返回不符合预期, 预期[%s], 实际[%s], 错误信息: %s", root.EDBEXECERROR, root.ErrorCode(err), err.Error())
	}

}

func TestUser_CreateSaltError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	crypto := Crypto{}
	crypto.saltFn = func(s string) (string, error) {
		return "", errors.New("加密错误")
	}
	session.mysqlSession.SetCrypto(&crypto)

	var user root.User
	user.No = "0001"
	user.Name = "测试"
	user.Password = "0001"

	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(user.No).WillReturnError(sql.ErrNoRows)

	service := session.mysqlSession.UserService()
	err = service.CreateUser(&user)
	if root.ErrorCode(err) != root.ECRYPTOERROR {
		t.Errorf("错误值返回不符合预期, 预期[%s], 实际[%s], 错误信息: %s", root.ECRYPTOERROR, root.ErrorCode(err), err.Error())
	}

}

func TestUser_LoginOK(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	credentials := root.Credentials{
		UserName: "0001",
		PassWord: "0001",
	}

	columns := []string{"no", "name", "password"}
	rows := sqlmock.NewRows(columns).
		AddRow("0001", "测试", "10001000")

	session.mock.ExpectQuery(`SELECT (.+) FROM users WHERE`).WithArgs(credentials.UserName).WillReturnRows(rows)

	service := session.mysqlSession.UserService()

	_, err = service.Login(credentials)
	if err != nil {
		t.Error(err)
	}

}

func TestUser_LoginQueryError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	credentials := root.Credentials{
		UserName: "0001",
		PassWord: "0001",
	}

	err = errors.New("查询出错")

	session.mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(credentials.UserName).WillReturnError(err)

	service := session.mysqlSession.UserService()

	_, err = service.Login(credentials)
	if root.ErrorCode(err) != root.EDBQUERYERROR {
		t.Errorf("错误返回值不符合预期，预期[%s],实际[%s],错误信息[%s]", root.EDBQUERYERROR, root.ErrorCode(err), err.Error())
	}

}

func TestUser_LoginWrongPWD(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	credentials := root.Credentials{
		UserName: "0001",
		PassWord: "0001",
	}

	columns := []string{"no", "password"}
	rows := sqlmock.NewRows(columns).
		AddRow("0001", "0002")

	session.mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(credentials.UserName).WillReturnRows(rows)

	service := session.mysqlSession.UserService()

	_, err = service.Login(credentials)
	if root.ErrorCode(err) != "password_invalid" {
		t.Errorf("错误返回类型不符合预期， 预期[%s]，实际[%s]，错误信息[%s]", "password_invalid", root.ErrorCode(err), err.Error())
	}

}

func TestUser_LoginCryptoError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer session.Close()

	var crypto Crypto
	crypto.compareFn = func(hash string, s string) (bool, error) {
		return false, errors.New("解加密错误")
	}

	session.mysqlSession.SetCrypto(&crypto)

	credentials := root.Credentials{
		UserName: "0001",
		PassWord: "0001",
	}

	columns := []string{"no", "password"}
	rows := sqlmock.NewRows(columns).
		AddRow("0001", "0002")

	session.mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(credentials.UserName).WillReturnRows(rows)

	service := session.mysqlSession.UserService()

	_, err = service.Login(credentials)
	if root.ErrorCode(err) != root.ECRYPTOERROR {
		t.Errorf("错误返回值不符合预期，预期[%s], 实际[%s], 错误信息:[%s]", root.ECRYPTOERROR, root.ErrorCode(err), err.Error())
	}

}

func TestUser_AuthError(t *testing.T) {

	session, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer session.Close()

	var auth Authenticator
	auth.tokenFn = func(*root.User) (string, error) {
		return "", errors.New("身份验证错误")
	}

	session.mysqlSession.SetAuthenticator(&auth)

	credentials := root.Credentials{
		UserName: "0001",
		PassWord: "0001",
	}

	columns := []string{"no", "password"}
	rows := sqlmock.NewRows(columns).
		AddRow("0001", "10001000")

	session.mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(credentials.UserName).WillReturnRows(rows)

	service := session.mysqlSession.UserService()

	_, err = service.Login(credentials)
	if root.ErrorCode(err) != "login_auth_error" {
		t.Errorf("错误返回值不符合预期，预期[%s], 实际[%s], 错误信息:[%s]", "login_auth_error", root.ErrorCode(err), err.Error())
	}

}
