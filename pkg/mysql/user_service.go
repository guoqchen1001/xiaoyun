package mysql

import (
	"database/sql"
	"errors"
	root "xiaoyun/pkg"
)

// UserService 用户服务
type UserService struct {
	session *Session
}

// Login 登录
func (s *UserService) Login(credentials root.Credentials) (string, error) {

	var customErr root.Error
	customErr.Op = "mysql.UserService.Login"

	userName := credentials.UserName
	passWord := credentials.PassWord

	var model UserNodel
	err := s.session.db.Get(&model, `SELECT	no, name ,password 
									 FROM users 
								 	 WHERE no = ? `, userName)
	if err != nil {
		customErr.Err = err
		customErr.Code = root.EDBQUERYERROR
		return "", &customErr
	}

	ok, err := s.session.crypto.Compare(model.Password, model.No+passWord)
	if err != nil {
		customErr.Code = root.ECRYPTOERROR
		customErr.Err = err
		return "", &customErr
	}

	if !ok {
		customErr.Code = "password_invalid"
		customErr.Err = errors.New("密码不正确")
		return "", &customErr
	}

	user := model.toUser()

	tokenStr, err := s.session.authenticator.Token(&user)
	if err != nil {
		customErr.Err = err
		customErr.Code = "login_auth_error"
		return "", &customErr
	}

	return tokenStr, nil

}

// CreateUser 创造用户
func (s *UserService) CreateUser(user *root.User) error {

	var customError root.Error
	customError.Op = "mysql.Session.CreateUser"

	model := toUserModel(user)

	passwrod, err := s.session.crypto.Salt(model.Password)
	if err != nil {
		customError.Err = err
		customError.Code = root.ECRYPTOERROR
		return &customError
	}
	model.Password = passwrod

	err = s.session.db.Get(&model, `SELECT no, name, password FROM users WHERE no = ?`, user.No)

	if err != nil {
		if err != sql.ErrNoRows {
			customError.Err = err
			customError.Code = root.EDBQUERYERROR
			return &customError
		}
	} else {
		customError.Code = root.ECONFLICT
		return &customError
	}

	// 密码加盐处理

	tx, err := s.session.db.Beginx()
	if err != nil {
		customError.Code = root.EDBBEGINERROR
		customError.Err = err
		return &customError
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamed(`INSERT INTO users (no, name, password)
								  VALUES(:no, :name, :password)`)
	if err != nil {
		customError.Code = root.EDBPREPAREERROR
		customError.Err = err
		return &customError
	}

	_, err = stmt.Exec(model)
	if err != nil {
		tx.Rollback()
		customError.Code = root.EDBEXECERROR
		customError.Err = err
		return &customError
	}

	tx.Commit()

	return nil
}
