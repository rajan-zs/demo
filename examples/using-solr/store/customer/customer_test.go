package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/examples/using-solr/store"
	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

const er = "error"

type test struct {
	collection string
	wantErr    bool
}

func TestCustomer_ListError(t *testing.T) {
	collections := []string{"error", "json error"}
	c := New(mockSolrClient{})
	k := gofr.New()
	ctx := gofr.NewContext(nil, nil, k)

	for _, collection := range collections {
		_, err := c.List(ctx, collection, store.Filter{})
		if err == nil {
			t.Error("Expected error got nil")
		}
	}
}

func TestCustomer_ListResponse(t *testing.T) {
	c := New(mockSolrClient{})
	k := gofr.New()
	ctx := gofr.NewContext(nil, nil, k)
	expectedResp := []store.Model{{ID: 553573403, Name: "book", DateOfBirth: "01-01-1987"}}

	resp, err := c.List(ctx, "customer", store.Filter{})
	if err != nil {
		t.Errorf("Expected nil error\tGot %v", err)
	}

	if !reflect.DeepEqual(resp, expectedResp) {
		t.Errorf("Expected %v\tGot %v\n", expectedResp, resp)
	}
}

func TestCustomer_Create(t *testing.T) {
	var testcases = []test{
		{"error", true},
		{"customer", false},
	}
	c := New(mockSolrClient{})
	k := gofr.New()
	ctx := gofr.NewContext(nil, nil, k)

	for _, tc := range testcases {
		err := c.Create(ctx, tc.collection, store.Model{})

		if (err == nil && tc.wantErr) || (err != nil && tc.wantErr == false) {
			t.Errorf("Expected %v\tGot %v\n", tc.wantErr, err)
		}
	}
}

func TestCustomer_Update(t *testing.T) {
	var testcases = []test{
		{"error", true},
		{"customer", false},
	}
	c := New(mockSolrClient{})
	k := gofr.New()
	ctx := gofr.NewContext(nil, nil, k)

	for _, tc := range testcases {
		err := c.Update(ctx, tc.collection, store.Model{})
		if (err == nil && tc.wantErr) || (err != nil && tc.wantErr == false) {
			t.Errorf("Expected %v\tGot %v\n", tc.wantErr, err)
		}
	}
}

func TestCustomer_Delete(t *testing.T) {
	var testcases = []test{
		{"error", true},
		{"customer", false},
	}
	c := New(mockSolrClient{})
	k := gofr.New()
	ctx := gofr.NewContext(nil, nil, k)

	for _, tc := range testcases {
		err := c.Delete(ctx, tc.collection, store.Model{})
		if (err == nil && tc.wantErr) || (err != nil && tc.wantErr == false) {
			t.Errorf("Expected %v\tGot %v\n", tc.wantErr, err)
		}
	}
}

type mockSolrClient struct{}

func (m mockSolrClient) Search(ctx context.Context, collection string, params map[string]interface{}) (interface{}, error) {
	if collection == er {
		return nil, errors.InvalidParam{}
	} else if collection == "json error" {
		b := []byte(`{"response": {
		"numFound": 1,
		"start": 0,
		"docs": [
			{	"id": "0553573403",
				"name": [
					"book"]}]}}`)
		var resp interface{}

		_ = json.Unmarshal(b, &resp)

		return datastore.Response{Code: 200, Data: resp}, nil
	}

	b := []byte(`{"response": {
		"numFound": 1,
		"start": 0,
		"docs": [
			{	"id": "553573403",
				"name":"book",
                "dateOfBirth":"01-01-1987"}]}}`)

	var resp interface{}
	_ = json.Unmarshal(b, &resp)

	return datastore.Response{Code: 200, Data: resp}, nil
}

func (m mockSolrClient) Create(c context.Context, collection string, d *bytes.Buffer, p map[string]interface{}) (interface{}, error) {
	if collection == er {
		return nil, errors.InvalidParam{}
	}

	b := []byte(`{"responseHeader": {
		"rf": 1,
    	"status": 0`)

	var resp interface{}

	_ = json.Unmarshal(b, &resp)

	return datastore.Response{Code: 200, Data: resp}, nil
}

func (m mockSolrClient) Update(c context.Context, collection string, d *bytes.Buffer, p map[string]interface{}) (interface{}, error) {
	if collection == "error" {
		return nil, errors.InvalidParam{}
	}

	b := []byte(`{"responseHeader": {
		"rf": 1,
    	"status": 0`)

	var resp interface{}
	_ = json.Unmarshal(b, &resp)

	return datastore.Response{Code: 200, Data: resp}, nil
}

func (m mockSolrClient) Delete(c context.Context, collection string, doc *bytes.Buffer, p map[string]interface{}) (interface{}, error) {
	if collection == "error" {
		return nil, errors.InvalidParam{}
	}

	b := []byte(`{"responseHeader": {
		"rf": 1,
    	"status": 0`)

	var resp interface{}
	_ = json.Unmarshal(b, &resp)

	return datastore.Response{Code: 200, Data: resp}, nil
}
