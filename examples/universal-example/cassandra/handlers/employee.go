package handlers

import (
	"strconv"

	"github.com/zopsmart/gofr/examples/universal-example/cassandra/entity"
	"github.com/zopsmart/gofr/examples/universal-example/cassandra/store"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type employee struct {
	model store.Employee
}

//nolint:golint //employee should not get accessed, without initialization of store.Employee
func New(e store.Employee) employee {
	return employee{
		model: e,
	}
}

func (e employee) Get(c *gofr.Context) (interface{}, error) {
	var filter entity.Employee

	params := c.Params()

	filter.ID, _ = strconv.Atoi(params["id"])
	filter.Name = params["name"]
	filter.Phone = params["phone"]
	filter.Email = params["email"]
	filter.City = params["city"]

	return e.model.Get(c, filter), nil
}

func (e employee) Create(c *gofr.Context) (interface{}, error) {
	var emp entity.Employee
	// json error
	if err := c.Bind(&emp); err != nil {
		return nil, err
	}

	records := e.model.Get(c, entity.Employee{ID: emp.ID})

	if len(records) > 0 {
		return nil, errors.EntityAlreadyExists{}
	}

	results, err := e.model.Create(c, emp)

	return results, err
}
