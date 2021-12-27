package gofr

import (
	"os"

	"go.opencensus.io/trace"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

type cmdApp struct {
	Router      CMDRouter
	metricSvr   *metricServer
	context     *Context
	tracingSpan *trace.Span
}

type metricServer struct {
	port  int
	route string
}

func (app *cmdApp) Start(logger log.Logger) {
	args := os.Args[1:] // 1st one is the command name itself.
	command := ""

	for _, a := range args {
		if a[1] != '-' {
			command = command + " " + a
		}
	}

	// start the server for health-check and metrics
	go func() {
		app := NewWithConfig(app.context.Config)
		app.Start()
	}()

	h := app.Router.handler(command)
	if h == nil {
		app.context.resp.Respond("No Command Found!", nil)
		return
	}

	data, err := h(app.context)
	if err != nil {
		app.context.resp.Respond(nil, err)
	} else {
		app.context.resp.Respond(data, nil)
	}

	app.tracingSpan.End()
}
