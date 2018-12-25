package main_test

import (
	"testing"
	main "xiaoyun/cmd/app"
)

type App struct {
	App *main.App
}

func TestApp_Init(t *testing.T) {
	var app App
	app.App = &main.App{}
	app.App.Initialize()
}
