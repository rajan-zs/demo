package main

import (
	"github.com/zopsmart/gofr/examples/sample-api/handler"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	// create the application object
	k := gofr.New()
	k.Server.ValidateHeaders = false
	// enabling /swagger endpoint for Swagger UI
	k.EnableSwaggerUI()

	// add a handler
	k.GET("/hello-world", handler.HelloWorld)

	// handler can access the parameters from context.
	k.GET("/hello", handler.HelloName)

	// handler function can send response in JSON using c.JSON
	k.GET("/json", handler.JSONHandler)

	// Handler function which throws error
	k.GET("/error", handler.ErrorHandler)

	// Handler function which uses logging
	k.GET("/log", handler.HelloLogHandler)

	// start the server
	k.Start()
}