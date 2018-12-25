package mssql

import (
	"time"
	root "xiaoyun/pkg"

	"github.com/jmoiron/sqlx"
)

// Session 数据库连接
type Session struct {
	db  *sqlx.DB
	now time.Time // 提供同一个会话中的数据库操作时间戳

	// 身份验证
	authenticator root.Authenticator
	authToken     string
	user          *root.User

	// 商品服务
	goodsService GoodsService
	// 外卖服务

}

// NewSession 创建链接
func NewSession(db *sqlx.DB) *Session {
	s := &Session{db: db}
	// 指定商品服务所用的数据库session
	s.goodsService.session = s
	return s
}

// SetAuthToken 设置认证token，实现xiaoyun.Seesion接口
func (s *Session) SetAuthToken(token string) {
	s.authToken = token
}

// Authenticate 通过token认证， 实现xiaoyun.Session接口
func (s *Session) Authenticate() (*root.User, error) {
	// 如果已通过身份验证，则直接返回用户
	if s.user != nil {
		return s.user, nil
	}

	// 通过session对象的认证接口进行身份认证
	u, err := s.authenticator.Authenticate(s.authToken)
	if err != nil {
		return nil, err
	}

	// 缓存已认证的用户信息
	s.user = u
	return u, nil

}

// GoodsService 商品服务接口，实现root.GoodsService
func (s *Session) GoodsService() root.GoodsService {
	return &s.goodsService
}

// SetAuthenticator 设置身份验证
func (s *Session) SetAuthenticator(a root.Authenticator) {
	s.authenticator = a
}
