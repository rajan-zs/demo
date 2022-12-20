package handler

import "developer.zopsmart.com/go/gofr/pkg/gofr"

type demo struct {
	Message string `json:"message"`
}

func Pong(c *gofr.Context) (interface{}, error) {
	var demo demo
	err := c.Bind(&demo)
	if err != nil {
		return nil, err
	}
	return demo, nil
}

func Ping(c *gofr.Context) (interface{}, error) {
	return "Hello from gofr ping", nil

}
