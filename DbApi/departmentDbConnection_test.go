package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	_ "net/http"
	"net/http/httptest"
	"testing"
)

var mock sqlmock.Sqlmock

func TestDbWithHttp(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	tests := []struct {
		description string
		reqBody     []byte
		expCode     int
		expResp     string
	}{
		{
			description: "Get All Departments ",
			reqBody:     nil,
			expCode:     200,
			expResp:     `[{"DepId":"49682d52-6ab5-44aa-a817-fa4c5b05e086","depName":"DEV"}]`,
		},
	}
	rows := sqlmock.NewRows([]string{"deptId", "depName"})
	rows.AddRow("49682d52-6ab5-44aa-a817-fa4c5b05e086", "DEV")
	mock.ExpectQuery("SELECT (.+) FROM dept").WillReturnRows(rows)
	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/departments", bytes.NewReader(tc.reqBody))
		mockResp := httptest.NewRecorder()
		deptHandlerGet(mockResp, mockReq)
		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)
		assert.Equal(t, tc.expResp, mockResp.Body.String())
	}

	mock.ExpectCommit()
}
func TestDeptHandlerPost(t *testing.T) {
	input := Department{
		"9682d52-6ab5-44aa-a817-fa4c5b05e086", "DEV",
	}
	reqBody, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/deptHandlerPost", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO dept").WithArgs("9682d52-6ab5-44aa-a817-fa4c5b05e086", "DEV").WillReturnResult(sqlmock.NewResult(1, 1))
	deptHandlerPost(w, req)
	resp := w.Result()
	if resp.StatusCode == http.StatusCreated {
		t.Log("Data is updated succesfully")
	}

}

func TestEmpHandlerGet(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	tests := []struct {
		description string
		reqBody     []byte
		expCode     int
		expResp     string
	}{
		{
			description: "Get All Employee ",
			reqBody:     nil,
			expCode:     200,
			expResp:     `[{"EmpId":"49682d52-6ab5-44aa-a817-fa4c5b05e086","EmpName":"Rajan","DeptId":"49682d52-6ab5-44aa-a817-fa4c5b05e086","Phone":1234,"depName":"DEV"}]`,
		},
	}
	rows := sqlmock.NewRows([]string{"EmpId", "EmpName", "DeptId", "Phone", "DepName"})
	rows.AddRow("49682d52-6ab5-44aa-a817-fa4c5b05e086", "Rajan", "49682d52-6ab5-44aa-a817-fa4c5b05e086", "1234", "DEV")
	mock.ExpectQuery("SELECT emp.empId,emp.empName,dept.deptId, emp.phone ,dept.depName FROM emp INNER JOIN dept on emp.id = dept.deptId").WillReturnRows(rows)
	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/employees", bytes.NewReader(tc.reqBody))

		mockResp := httptest.NewRecorder()
		empHandlerGet(mockResp, mockReq)
		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)
		assert.Equal(t, tc.expResp, mockResp.Body.String())
	}

	mock.ExpectCommit()
}

func TestEmpHandlerPost(t *testing.T) {
	input := Employee{
		"49682d52-6ab5-44aa-a817-fa4c5b05e086", "Rajan", "49682d52-6ab5-44aa-a817-fa4c5b05e086", 1234, "DEV",
	}
	db, mock, err = sqlmock.New()
	reqBody, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/employeePost", bytes.NewBuffer(reqBody))

	w := httptest.NewRecorder()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO emp").WithArgs("49682d52-6ab5-44aa-a817-fa4c5b05e086", "Rajan", "49682d52-6ab5-44aa-a817-fa4c5b05e086", 1234, "DEV").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	empHandlerPost(w, req)
	resp := w.Result()
	if resp.StatusCode == http.StatusCreated {
		t.Log("Data is updated succesfully")
	}

}
func TestGetEmpByIdHandler(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	tests := []struct {
		description string
		reqBody     []byte
		expCode     int
		expResp     string
	}{
		{
			description: "Get All Employee ",
			reqBody:     nil,
			expCode:     200,
			expResp:     `[{"EmpId":"49682d52-6ab5-44aa-a817-fa4c5b05e086","EmpName":"Rajan","DeptId":"49682d52-6ab5-44aa-a817-fa4c5b05e086","Phone":1234,"depName":"DEV"}]`,
		},
	}

	rows := sqlmock.NewRows([]string{"EmpId", "EmpName", "DeptId", "Phone", "DepName"})
	rows.AddRow("49682d52-6ab5-44aa-a817-fa4c5b05e086", "Rajan", "49682d52-6ab5-44aa-a817-fa4c5b05e086", "1234", "DEV")
	mock.ExpectQuery("SELECT emp.empId, emp.empName, dept.deptId, emp.phone ,dept.depName FROM emp INNER JOIN dept on emp.id = dept.deptId ").WillReturnRows(rows)
	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/employees?=49682d52-6ab5-44aa-a817-fa4c5b05e086", bytes.NewReader(tc.reqBody))
		mockResp := httptest.NewRecorder()
		getEmpByIdHandler(mockResp, mockReq)
		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)
		assert.Equal(t, tc.expResp, mockResp.Body.String())
	}

	mock.ExpectCommit()
}
