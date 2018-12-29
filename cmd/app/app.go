package main

import (
	root "xiaoyun/pkg"
	"xiaoyun/pkg/crypto"
	"xiaoyun/pkg/http"
	"xiaoyun/pkg/jwt"
	"xiaoyun/pkg/mssql"
	"xiaoyun/pkg/mysql"
)

//App 主程序
type App struct {
	env         *Env
	path        *root.PathConfig   // 路径配置
	configer    root.Configer      // app配置对象
	log         *root.Log          // 日志对象
	auth        root.Authenticator // 身份验证
	clientMssql *mssql.Client      // mssql连接对象
	clientMysql *mysql.Client      // mysql连接对象
	handler     *http.Handler      // http处理器
	sever       *http.Server       // http服务
	crypto      *crypto.Crypto     // 加密服务

	sessionMssql *mssql.Session
	sessionMysql *mysql.Session
}

// SetEnv 初始化系统运行环境
func (a *App) SetEnv(env *Env) {
	a.configer = env.Configer
	a.log = env.Log
}

// Initialize 初始化
func (a *App) Initialize() {

	a.initializeCrypto()
	a.initializeAuth()
	a.initializeMssql()
	a.initializeMysql()
	a.initializeHanlder()
	a.initializeServer()
}

//Run 启动程序
func (a *App) Run() error {

	// 设置运行环境
	if err := a.clientMssql.Open(); err != nil {
		return err
	}
	defer a.clientMssql.Close()
	if err := a.clientMysql.Open(); err != nil {
		return err
	}

	err := a.clientMysql.Migrate(a.log)
	if err != nil {
		return err
	}

	defer a.clientMysql.Close()

	a.initializeService()

	return a.sever.Open()
}

// 初始化认证对象
func (a *App) initializeAuth() {
	authenticator := jwt.NewAuthenticator(a.configer)
	a.auth = authenticator

}

// 初始化加密服务
func (a *App) initializeCrypto() {
	crypto := crypto.NewCrypto()
	a.crypto = crypto
}

// 初始化mssql对象
func (a *App) initializeMssql() {
	// mssql数据库单次链接
	client := mssql.NewClient(a.configer)
	a.clientMssql = client
}

// 初始化mysql对象
func (a *App) initializeMysql() {
	client := mysql.NewClient(a.configer)
	client.Authenticator = a.auth
	client.Crypto = a.crypto.MD5Crypto()
	a.clientMysql = client
}

// 初始化http处理器对象
func (a *App) initializeHanlder() {
	handler := http.NewHandler(a.log)
	a.handler = handler
}

// 初始化http server对象
func (a *App) initializeServer() {
	server := http.NewServer(a.configer, a.handler)
	a.sever = server
}

func (a *App) initializeService() {

	a.sessionMssql = a.clientMssql.Connect()
	a.sessionMysql = a.clientMysql.Connect()

	a.handler.GoodsHandler.GoodsService = a.sessionMssql.GoodsService()
	a.handler.GoodsImageHandler.GoodsImageService = a.sessionMysql.GoodsImageService()
	a.handler.UserHandler.UserService = a.sessionMysql.UserService()

}
