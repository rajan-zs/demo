package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zopsmart/gofr/examples/universal-example/pgsql/entity"
	gofrError "github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/gofr/responder"
)

type mockStore struct{}

const (
	fetchErr  = constError("error while fetching employee listing")
	createErr = constError("error while adding new employee")
)

type constError string

func (err constError) Error() string {
	return string(err)
}

func (m mockStore) Get(c *gofr.Context) ([]entity.Employee, error) {
	p := c.Param("mock")
	if p == "success" {
		return nil, nil
	}

	return nil, fetchErr
}

func (m mockStore) Create(c *gofr.Context, customer entity.Employee) error {
	switch customer.Name {
	case "some_employee":
		return nil
	case "mock body error":
		return gofrError.InvalidParam{Param: []string{"body"}}
	}

	return createErr
}

func TestPgsqlEmployee_Get(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		mockParamStr string
		expectedErr  error
	}{
		{"mock=success", nil},
		{"", fetchErr},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/dummy?"+tc.mockParamStr, nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		c := gofr.NewContext(res, req, k)

		_, err := m.Get(c)
		assert.Equal(t, tc.expectedErr, err, i)
	}
}

func TestPgsqlEmployee_Create(t *testing.T) {
	m := New(mockStore{})
	k := gofr.New()

	tests := []struct {
		body        []byte
		expectedErr error
	}{
		{[]byte(`{"name":"some_employee"}`), nil},
		{[]byte(`mock body error`), gofrError.InvalidParam{Param: []string{"body"}}},
		{[]byte(`{"name":"creation error"}`), createErr},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://dummy", bytes.NewReader(tc.body))
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		c := gofr.NewContext(res, req, k)

		_, err := m.Create(c)
		assert.Equal(t, tc.expectedErr, err, i)
	}
}