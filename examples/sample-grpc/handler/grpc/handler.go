package grpc

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type handler struct {
	UnimplementedExampleServiceServer
}

// New is factory function for GRPC Handler
//nolint:revive // handler should not be used without proper initilization with required dependency
func New() handler {
	return handler{}
}

func (h handler) Get(ctx context.Context, id *ID) (*Response, error) {
	if id.Id == "1" {
		resp := &Response{
			FirstName:  "First",
			SecondName: "Second",
		}

		return resp, nil
	}

	return nil, errors.EntityNotFound{Entity: "name", ID: id.Id}
}
