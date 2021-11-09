package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type store struct{}

// New is factory function for customer
//nolint:revive // customer should not be used without proper initilization with required dependency
func New() store {
	return store{}
}

const index = "customers"

func (s store) Get(ctx *gofr.Context, name string) ([]model.Customer, error) {
	var body string

	if name != "" {
		body = fmt.Sprintf(`{"query" : { "match" : {"name":"%q"} }}`, name)
	}

	es := ctx.Elasticsearch

	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithContext(ctx),
		es.Search.WithBody(strings.NewReader(body)),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	var customers []model.Customer

	err = es.BindArray(res, &customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s store) GetByID(ctx *gofr.Context, id string) (model.Customer, error) {
	var customer model.Customer

	es := ctx.Elasticsearch

	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithContext(ctx),
		es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{"query" : { "match" : {"id":"%q"} }}`, id))),
		es.Search.WithPretty(),
		es.Search.WithSize(1),
	)
	if err != nil {
		return customer, errors.DB{Err: err}
	}

	err = es.Bind(res, &customer)
	if err != nil {
		return customer, err
	}

	if customer.ID == "" {
		return customer, errors.EntityNotFound{Entity: "customer", ID: id}
	}

	return customer, nil
}

func (s store) Update(ctx *gofr.Context, c model.Customer, id string) (model.Customer, error) {
	body, err := json.Marshal(c)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	es := ctx.Elasticsearch

	res, err := es.Index(
		index,
		bytes.NewReader(body),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithContext(ctx),
		es.Index.WithDocumentID(id),
	)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	resp, err := es.Body(res)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	if id, ok := resp["_id"].(string); ok {
		return s.GetByID(ctx, id)
	}

	return model.Customer{}, errors.Error("update error: invalid id")
}

func (s store) Create(ctx *gofr.Context, c model.Customer) (model.Customer, error) {
	body, err := json.Marshal(c)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	es := ctx.Elasticsearch

	res, err := es.Index(
		index,
		bytes.NewReader(body),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithContext(ctx),
		es.Index.WithDocumentID(c.ID),
	)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	resp, err := es.Body(res)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	if id, ok := resp["_id"].(string); ok {
		return s.GetByID(ctx, id)
	}

	return model.Customer{}, errors.Error("create error: invalid id")
}

func (s store) Delete(ctx *gofr.Context, id string) error {
	es := ctx.Elasticsearch

	_, err := es.Delete(
		index,
		id,
		es.Delete.WithContext(ctx),
		es.Delete.WithPretty(),
	)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
