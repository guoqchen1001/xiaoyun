package main

func main() {

	a := App{}

	// 初始化

	env := NewProductEnv()
	a.SetEnv(env)

	a.Initialize()
	// 运行
	err := a.Run()
	if err != nil {
		a.log.Logger.Error(err)
	}

}
