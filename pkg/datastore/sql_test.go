package datastore

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func getDB() DataStore {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")

	dc := DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	}

	db, _ := NewORM(&dc)

	store := new(DataStore)

	store.rdb.DB = db.DB.DB()
	store.rdb.config = db.config
	store.rdb.logger = log.NewMockLogger(io.Discard)

	return *store
}

func TestSQLClient_Exec(t *testing.T) {
	db := getDB()

	_, err := db.DB().Exec("SHOW TABLES")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}
}

func TestSQLClient_Query(t *testing.T) {
	db := getDB()

	_, err := db.DB().Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	rows, err := db.DB().Query("SELECT * FROM testtable")
	if err != nil {
		t.Errorf("Query operation failed. Got: %s", err)
	}

	if rows.Err() != nil {
		t.Errorf("Encountered error: %s", rows.Err())
	}

	defer rows.Close()

	if rows == nil {
		t.Errorf("Failed. Got empty rows")
	}
}

func TestSQLClient_QueryRow(t *testing.T) {
	db := getDB()

	_, err := db.DB().Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	row := db.DB().QueryRow("SELECT * FROM testtable")
	if row == nil {
		t.Errorf("QueryRow operation failed.")
	}
}

func TestSQLClient_QueryContext(t *testing.T) {
	db := getDB()

	_, err := db.DB().Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	rows, err := db.DB().QueryContext(context.Background(), "SELECT * FROM testtable")
	if err != nil {
		t.Errorf("Query operation failed. Got: %s", err)
	}

	if rows.Err() != nil {
		t.Errorf("Encountered error: %s", rows.Err())
	}

	defer rows.Close()

	if rows == nil {
		t.Errorf("Failed. Got empty rows")
	}
}

func TestSQLClient_ExecContext(t *testing.T) {
	db := getDB()

	_, err := db.DB().ExecContext(context.Background(), "SHOW TABLES")
	if err != nil {
		t.Errorf("ExecContext operation failed. Got: %s", err)
	}
}

func TestSQLClient_QueryRowContext(t *testing.T) {
	db := getDB()

	_, err := db.DB().Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = db.DB().Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	row := db.DB().QueryRowContext(context.Background(), "SELECT * FROM testtable")
	if row == nil {
		t.Errorf("QueryRow operation failed.")
	}
}

func Test_checkQueryOperation(t *testing.T) {
	queries := []string{"SELECT * FROM randomTable", "DELETE FROM randomTable", "UPDATE randomTable",
		"INSERT INTO randomTable", "insert INTO (SELECT *)", "    SELECT * FROM", "\n  UPDATE TABLE",
		"\nDELETE FROM", "BEGIN TR", "COMMIT TR", "SET <EXPR>"}

	expected := []string{"SELECT", "DELETE", "UPDATE", "INSERT", "INSERT", "SELECT", "UPDATE", "DELETE", "BEGIN", "COMMIT", "SET"}

	for i := range queries {
		operation := checkQueryOperation(queries[i])
		if operation != expected[i] {
			t.Errorf("Testcase %v Failed. Expected: %s, Got: %s", i, expected[i], operation)
		}
	}
}

func Test_QueryLog(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	c := config.NewGoDotEnvProvider(logger, "../../configs")

	dc := DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	}

	db, _ := NewORM(&dc)

	ds := new(DataStore)

	ds.rdb.DB = db.DB.DB()
	ds.rdb.config = db.config
	ds.rdb.logger = logger

	query := "SELECT NOW()"

	_, _ = ds.DB().Exec(query)

	if !strings.Contains(b.String(), query) {
		t.Errorf("Expected %v\nGot %v", query, b.String())
	}

	if !strings.Contains(b.String(), "sql") {
		t.Errorf("Expected %vGot %v", "SQL", b.String())
	}
}

func TestSQLTx_Exec(t *testing.T) {
	db := getDB()

	tx, err := db.DB().Begin()
	if err != nil {
		t.Errorf("Begin transaction operation failed. Got: %s", err)
	}

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}
}

func TestSQLTx_ExecContext(t *testing.T) {
	db := getDB()

	tx, err := db.DB().Begin()
	if err != nil {
		t.Errorf("Begin operation failed. Got: %s", err)
	}

	_, err = tx.ExecContext(context.Background(), "SHOW TABLES")
	if err != nil {
		t.Errorf("ExecContext operation failed. Got: %s", err)
	}
}

func TestSQLTx_Query(t *testing.T) {
	db := getDB()

	tx, err := db.DB().Begin()
	if err != nil {
		t.Errorf("Begin operation failed. Got: %s", err)
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	rows, err := tx.Query("SELECT * FROM testtable")
	if err != nil {
		t.Errorf("Query operation failed. Got: %s", err)
	}

	if rows.Err() != nil {
		t.Errorf("Encountered error: %s", rows.Err())
	}

	defer rows.Close()

	if rows == nil {
		t.Errorf("Failed. Got empty rows")
	}
}

func TestSQLTx_QueryContext(t *testing.T) {
	db := getDB()

	tx, err := db.DB().Begin()
	if err != nil {
		t.Errorf("Begin operation failed. Got: %s", err)
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	rows, err := tx.QueryContext(context.Background(), "SELECT * FROM testtable")
	if err != nil {
		t.Errorf("Query operation failed. Got: %s", err)
	}

	if rows.Err() != nil {
		t.Errorf("Encountered error: %s", rows.Err())
	}

	defer rows.Close()

	if rows == nil {
		t.Errorf("Failed. Got empty rows")
	}
}

func TestSQLTx_QueryRow(t *testing.T) {
	db := getDB()

	tx, err := db.DB().Begin()
	if err != nil {
		t.Errorf("Begin operation failed. Got: %s", err)
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	row := tx.QueryRow("SELECT * FROM testtable")
	if row == nil {
		t.Errorf("QueryRow operation failed.")
	}
}

func TestSQLTx_QueryRowContext(t *testing.T) {
	db := getDB()

	tx, err := db.DB().Begin()
	if err != nil {
		t.Errorf("Begin operation failed. Got: %s", err)
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS testtable")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS testtable(id int PRIMARY KEY NOT NULL)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	_, err = tx.Exec("INSERT INTO testtable(id) values(32)")
	if err != nil {
		t.Errorf("Exec operation failed. Got: %s", err)
	}

	row := tx.QueryRowContext(context.Background(), "SELECT * FROM testtable")
	if row == nil {
		t.Errorf("QueryRow operation failed.")
	}
}

func TestSQLClient_BeginTx(t *testing.T) {
	db := getDB()

	_, err := db.DB().BeginTx(context.Background(), nil)
	if err != nil {
		t.Errorf("Error in starting the transaction")
	}
}

func TestSQLTx_Commit(t *testing.T) {
	db := getDB()

	tx, err := db.DB().BeginTx(context.Background(), nil)
	if err != nil {
		t.Errorf("Error in starting the transaction")
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf("Error encountered while committing the transaction")
	}
}

func Test_DataBaseNameInTransaction(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	dc := DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	}

	db, _ := NewORM(&dc)

	store := new(DataStore)

	store.rdb.DB = db.DB.DB()
	store.rdb.config = db.config
	store.rdb.logger = log.NewMockLogger(b)

	txn, _ := store.DB().Begin()

	if !strings.Contains(b.String(), "sql") {
		t.Errorf("Failed.\tExpected %v\tGot %v\n", "sql", b.String())
	}

	b.Reset()

	_, _ = txn.Exec("SHOW TABLES")

	if !strings.Contains(b.String(), "sql") {
		t.Errorf("Failed.\tExpected %v\tGot %v\n", "sql", b.String())
	}

	b.Reset()

	_ = txn.Commit()

	if !strings.Contains(b.String(), "sql") {
		t.Errorf("Failed.\tExpected %v\tGot %v\n", "sql", b.String())
	}
}
