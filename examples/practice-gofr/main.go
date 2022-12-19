package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()
	app.GET("/ping", ping)
	app.Start()
}

func ping(c *gofr.Context) (interface{}, error) {
	return "Hello from gofr ping", nil

}
