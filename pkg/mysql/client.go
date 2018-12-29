package mysql

import (
	"time"
	root "xiaoyun/pkg"

	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jmoiron/sqlx"
)

// Client 数据库连接
type Client struct {
	db            *sqlx.DB
	Now           func() time.Time
	Configer      root.Configer
	Authenticator root.Authenticator
	Crypto        root.Crypto
}

// NewClient 生成数据路客户端
func NewClient(configer root.Configer) *Client {

	c := Client{Configer: configer}
	c.Now = time.Now
	return &c

}

// Connect 返回数据库连接Connect对象,并打开数据库连接
func (c *Client) Connect() *Session {

	s := NewSession(c.db)
	s.SetAuthenticator(c.Authenticator)
	s.SetCrypto(c.Crypto)
	s.now = c.Now()
	return s

}

// Open 连接数据库
func (c *Client) Open() error {

	const op = "mysql.Client.Open"

	customErr := root.Error{
		Op: op,
	}

	config, err := c.Configer.GetConfig()
	if err != nil {
		customErr.Err = err
		customErr.Code = root.ECONFIGNOTINVALID
		return &customErr
	}

	if config.Mysql == nil {
		customErr.Code = root.ECONFIGMYSQLNOTFOUND
		return &customErr
	}

	dialects := config.Mysql.Dialects
	if dialects == "" {
		dialects = "mysql"
	}

	db, err := sqlx.Open(dialects, config.Mysql.Parm)
	if err != nil {
		customErr.Code = root.EDBOPENERROR
		customErr.Err = err
		return &customErr
	}

	c.db = db

	return nil

}

// Close 关闭数据库连接
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Migrate 数据库迁移
func (c *Client) Migrate(log *root.Log) error {

	config, err := c.Configer.GetConfig()
	if err != nil {
		return err
	}

	m, err := NewMigrate(c.db.DB, config.Mysql.Migrations, log.Logger)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
