package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const ENDPOINT = "http://localhost:8080/employee" // ENDPOINT constant for the endpoint URL

// Test_EmployeeHandler
func Test_EmployeeHandler(t *testing.T) {
	type testcase struct {
		input       Employee
		expected    interface{}
		statusCode  int
		method      string
		description string
	}

	// Test Cases
	var testcases = []testcase{
		{
			Employee{},
			ResponseMessageStruct{"Data is not avilable"},
			http.StatusNoContent,
			http.MethodGet,
			"Data is not avilable",
		},
		{
			Employee{
				"INT195",
				"Shashank Shekhar",
				"Patna",
				24,
			},
			ResponseMessageStruct{"successfully"},
			http.StatusCreated,
			http.MethodPost,
			"Post request -Data should be added because",
		},
		{
			Employee{},
			[]Employee{
				{"INT191", "Rajan", "Patna", 23},
			},
			http.StatusOK,
			http.MethodGet,
			"GET request",
		},
		{
			Employee{
				"INT196",
				"Mohit Bajaj",
				"Purnea",
				24,
			},
			ResponseMessageStruct{"successfully"},
			http.StatusCreated,
			http.MethodPost,
			"Data should be retrieved because it's a GET request",
		},
		{
			Employee{},
			[]Employee{
				{"INT195", "Shashank Shekhar", "Patna", 24},
				{"INT196", "Mohit Bajaj", "Purnea", 24},
			},
			http.StatusOK,
			http.MethodGet,
			"Data should be retrieved because it's a GET request",
		},
		{
			Employee{},
			ResponseMessageStruct{"Method not allowed"},
			http.StatusMethodNotAllowed,
			http.MethodPut,
			"Only POST and GET is allowed",
		},
	}

	var response ResponseMessageStruct
	var employees []Employee

	for i := range testcases {
		handler := EmployeeHandler
		input, _ := json.Marshal(testcases[i].input)
		req := httptest.NewRequest(testcases[i].method, ENDPOINT, bytes.NewBuffer([]byte(input)))
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
			json.Unmarshal(body, &response)
			if !reflect.DeepEqual(response, testcases[i].expected) {
				t.Errorf("Testcase: %d\n"+
					"Input: %v\n"+
					"Expected Output: %v\n"+
					"Actual Output: %v\n"+
					"Description: %v", i+1, testcases[i].input, testcases[i].expected, response, testcases[i].description)
			}
		}
	}
}
