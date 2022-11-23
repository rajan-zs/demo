package Assignement1Mux

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetRequest(t *testing.T) {
	type testcase struct {
		input       Employee
		expected    interface{}
		statusCode  int
		method      string
		description string
	}

	var testcases = []testcase{
		{
			Employee{},
			[]Employee{},
			http.StatusNoContent,
			http.MethodGet,
			"data is not avilable",
		},
		{
			Employee{},
			[]Employee{
				{"INT195", "Rajan", 21},
			},
			http.StatusOK,
			http.MethodGet,
			"Actal data is not matched with expected data",
		},
	}

	var employees = make([]Employee, 0)

	for i := range testcases {
		handler := handlerGet
		input, _ := json.Marshal(testcases[i].input)
		req := httptest.NewRequest(testcases[i].method, "http://localhost:8080/employee", bytes.NewBuffer([]byte(input)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		switch testcases[i].method {
		case http.MethodGet:
			json.Unmarshal(body, &employees)
			if !reflect.DeepEqual(employees, testcases[i].expected) {
				t.Errorf("Testcase: %d\n"+
					"Input: %v\n"+
					"Expected Output: %v\n"+
					"Actual Output: %v\n"+
					"Description: %v", i+1, testcases[i].input, testcases[i].expected, employees, testcases[i].description)
			}
		default:
			if testcases[i].expected == resp {
				t.Errorf("Testcase: %d\n"+
					"Input: %v\n"+
					"Expected Output: %v\n"+
					"Description: %v", i+1, testcases[i].input, testcases[i].expected, testcases[i].description)
			}
		}
	}
}
func TestPostRequest(t *testing.T) {
	type testcase struct {
		input       Employee
		expected    interface{}
		statusCode  int
		method      string
		description string
	}

	var testcases = []testcase{

		{
			Employee{
				"143",
				"RAjan",
				34,
			},
			"Data Added successfully",
			http.StatusCreated,
			http.MethodPost,
			"POST Request -Data will be added",
		},

		{
			Employee{
				"INT196",
				"Mohit Bajaj",
				24,
			},
			"Data Added successfully",
			http.StatusCreated,
			http.MethodPost,
			"POST Request -Data will be added",
		},

		{
			Employee{},
			"Method not allowed",
			http.StatusMethodNotAllowed,
			http.MethodPut,
			"Only POST and GET is allowed",
		},
	}

	var employees []Employee
	handler := handlerPost
	for i := range testcases {
		input, _ := json.Marshal(testcases[i].input)
		req := httptest.NewRequest(testcases[i].method, "http://localhost:8080/employee", bytes.NewBuffer([]byte(input)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		switch testcases[i].method {
		case http.MethodPost:
			json.Unmarshal(body, &employees)
			if !reflect.DeepEqual(employees, testcases[i].expected) {
				t.Errorf("Testcase: %d\n"+
					"Input: %v\n"+
					"Expected Output: %v\n"+
					"Actual Output: %v\n"+
					"Description: %v", i+1, testcases[i].input, testcases[i].expected, employees, testcases[i].description)
			}
		default:
			if testcases[i].expected == resp {
				t.Errorf("Testcase: %d\n"+
					"Input: %v\n"+
					"Expected Output: %v\n"+
					"Description: %v", i+1, testcases[i].input, testcases[i].expected, testcases[i].description)
			}
		}
	}
}
