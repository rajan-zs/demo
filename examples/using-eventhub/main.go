package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-eventhub/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	k.GET("/pub", handlers.Producer)
	k.GET("/sub", handlers.Consumer)

	k.Server.HTTP.Port = 9113
	k.Start()
}
