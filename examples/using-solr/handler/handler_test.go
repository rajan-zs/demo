package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/examples/using-solr/store"
	errors2 "github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
)

const er = "error"

func TestCustomer_List(t *testing.T) {
	testcases := []struct {
		query string
		err   error
	}{
		{"id=1&name=Henry", nil},
		{"id=123&name=Tomato", errors.New("core error")},
		{"", errors2.MissingParam{Param: []string{"id"}}},
	}
	c := New(&mockStore{})
	k := gofr.New()

	for i, tc := range testcases {
		req := httptest.NewRequest("GET", "/dummy?"+tc.query, nil)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		_, err := c.List(context)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("[TEST ID %d]Expected %v\tGot %v\n", i+1, tc.err, err)
		}
	}
}

func TestCustomer_Create(t *testing.T) {
	//nolint: govet, table tests
	testcases := []struct {
		body []byte
		err  error
	}{
		{[]byte(`{"id":1,"name":"Ethen"}`), nil},
		{[]byte(`{"id":1,"name":"error"}`), errors.New("core error")},

		{[]byte(`{"id":1}`), errors2.InvalidParam{[]string{"name"}}},

		{[]byte(`{"id":"1"}`), errors2.InvalidParam{[]string{"body"}}},
	}

	c := New(&mockStore{})
	k := gofr.New()

	for i, tc := range testcases {
		req := httptest.NewRequest("POST", "/dummy", bytes.NewBuffer(tc.body))
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		_, err := c.Create(context)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("[TEST CASE %d]Expected %v\tGot %v\n", i+1, tc.err, err)
		}
	}
}

func TestCustomer_Update(t *testing.T) {
	//nolint: govet, table tests
	testcases := []struct {
		body []byte
		err  error
	}{
		{[]byte(`{"id":1,"name":"Ethen"}`), nil},
		{[]byte(`{"id":1,"name":"error"}`), errors.New("core error")},
		{[]byte(`{"id":1}`), errors2.InvalidParam{Param: []string{"name"}}},
		{[]byte(`{"id":"1"}`), errors2.InvalidParam{[]string{"body"}}},
		{[]byte(`{"name":"Wen"}`), errors2.InvalidParam{[]string{"id"}}},
	}

	c := New(&mockStore{})
	k := gofr.New()

	for i, tc := range testcases {
		req := httptest.NewRequest("PUT", "/dummy", bytes.NewBuffer(tc.body))
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		_, err := c.Update(context)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("[TEST CASE %d]Expected %v\tGot %v\n", i+1, tc.err, err)
		}
	}
}

func TestCustomer_Delete(t *testing.T) {
	testcases := []struct {
		body []byte
		err  error
	}{
		{[]byte(`{"id":1,"name":"Ethen"}`), nil},
		{[]byte(`{"id":1,"name":"error"}`), errors.New("core error")},
		{[]byte(`{"id":"1"}`), errors2.InvalidParam{Param: []string{"body"}}},
		{[]byte(`{"name":"Wen"}`), errors2.InvalidParam{Param: []string{"id"}}},
	}

	c := New(&mockStore{})
	k := gofr.New()

	for i, tc := range testcases {
		req := httptest.NewRequest("DELETE", "/dummy", bytes.NewBuffer(tc.body))
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		_, err := c.Delete(context)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("[TEST CASE %d]Expected %v\tGot %v\n", i+1, tc.err, err)
		}
	}
}

func TestCustomer_Create2(t *testing.T) {
	c := New(&mockStore{})
	k := gofr.New()
	req := httptest.NewRequest("POST", "/dummy", errReader(0))
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)

	_, err := c.Delete(context)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}

	_, err = c.Create(context)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}

	_, err = c.Update(context)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

type mockStore struct{}

func (m *mockStore) List(ctx *gofr.Context, collection string, filter store.Filter) ([]store.Model, error) {
	if filter.ID == "1" {
		return []store.Model{{ID: 1, Name: "Henry", DateOfBirth: "01-01-1987"}}, nil
	}

	return nil, errors.New("core error")
}

func (m *mockStore) Create(ctx *gofr.Context, collection string, customer store.Model) error {
	if customer.Name == er {
		return errors.New("core error")
	}

	return nil
}

func (m *mockStore) Update(ctx *gofr.Context, collection string, customer store.Model) error {
	if customer.Name == "error" {
		return errors.New("core error")
	}

	return nil
}

func (m *mockStore) Delete(ctx *gofr.Context, collection string, customer store.Model) error {
	if customer.Name == "error" {
		return errors.New("core error")
	}

	return nil
}