package main

import (
	root "xiaoyun/pkg"
	"xiaoyun/pkg/config"
)

// Env 云鼎环境
type Env struct {
	No       string
	Name     string
	Log      *root.Log
	Configer root.Configer
}

// NewProductEnv 返回生产环境
func NewProductEnv() *Env {
	env := Env{
		No:   "product",
		Name: "生产环境",
	}

	env.Log = root.NewLogMultiOut("app.log")
	env.Configer = config.NewConfig("app_config.json")

	return &env
}

// NewDevelopEnv 返回开发环境
func NewDevelopEnv() *Env {
	env := Env{
		No:   "develop",
		Name: "开发环境",
	}

	env.Log = root.NewLogStdOut()

	return &env
}
