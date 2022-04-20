package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestMain(m *testing.M) {
	app := gofr.New()

	db := app.DB()
	if db == nil {
		return
	}

	query := `
 	   CREATE TABLE IF NOT EXISTS customers (
	   id serial primary key,
	   name varchar (50))
	`

	if app.Config.Get("DB_DIALECT") == "mssql" {
		query = `
		IF NOT EXISTS
	(  SELECT [name]
		FROM sys.tables
      WHERE [name] = 'customers'
   ) CREATE TABLE customers (id int primary key identity(1,1),
	   name varchar (50))
	`
	}

	if _, err := db.Exec(query); err != nil {
		app.Logger.Errorf("got error sourcing the schema: ", err)
	}

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(time.Second * 5)

	req, _ := request.NewMock(http.MethodGet, "http://localhost:9092/customer", nil)
	c := http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		t.Errorf("TEST Failed.\tHTTP request encountered Err: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed.\tExpected %v\tGot %v\n", http.StatusOK, resp.StatusCode)
	}

	_ = resp.Body.Close()
}