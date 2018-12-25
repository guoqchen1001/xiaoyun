package mysql

import (
	"time"
	root "xiaoyun/pkg"

	"github.com/jmoiron/sqlx"
)

// Session 数据库连接
type Session struct {
	db  *sqlx.DB
	now time.Time

	Authenticator root.Authenticator
	authToken     string
	user          *root.User

	goodsImageService GoodsImageService
}

// NewSession 创建数据库连接
func NewSession(db *sqlx.DB) *Session {
	s := &Session{db: db}

	s.goodsImageService.session = s

	return s
}

// SetAuthToken 设置token
func (s *Session) SetAuthToken(token string) {
	s.authToken = token
}

// Authenticate 身份验证，实现root.Authenticator接口
func (s *Session) Authenticate(token string) (*root.User, error) {

	customErr := root.Error{
		Op: "mysql.Session.Authenticate",
	}
	if s.user != nil {
		return s.user, nil
	}

	user, err := s.Authenticator.Authenticate(token)
	if err != nil {
		customErr.Code = root.EAUTHERROR
		customErr.Err = err
		return nil, &customErr
	}

	s.user = user

	return user, nil

}

// GoodsImageService 实现商品保存接口
func (s *Session) GoodsImageService() root.GoodsImageService {
	return &s.goodsImageService
}
