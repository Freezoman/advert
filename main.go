package main

import (
	"test/app"
	"test/model"
	"test/server"
)

func main() {
	m := model.Init()
	a := app.Init(m)
	s := server.Init(a)
	s.ServerStart()
}
