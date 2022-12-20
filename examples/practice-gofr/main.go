package main

import (
	"developer.zopsmart.com/go/gofr/examples/practice-gofr/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()
	app.Server.ValidateHeaders = false
	app.GET("/ping", handler.Ping)
	app.POST("/pong", handler.Pong)
	app.Start()
}
