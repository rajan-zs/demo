package main

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	_ "net/http"
	"testing"
)

type testcase struct {
	description        string
	ExpectedstatusCode int
}

var mock sqlmock.Sqlmock

func TestDbWithHttp(t *testing.T) {
	//db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"deptId", "depName"})
	rows.AddRow("49682d52-6ab5-44aa-a817-fa4c5b05e086 ", "DEV")
	mock.ExpectQuery("SELECT * from dept").WillReturnRows(rows)

	mock.ExpectCommit()

	handler := deptHandlerGet
	tcase := []testcase{{"cheking for success case", 405},
		{"checking for get request ", 500},
	}
}
