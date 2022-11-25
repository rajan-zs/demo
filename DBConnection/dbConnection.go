package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type employee struct {
	ID     string
	Name   string
	deptId string
	phone  string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "rajan:mypassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)

		return
	}

	employees, err := getEmployees(db)
	if err != nil {
		log.Println(err)
		return
	}

	oneEmployee, err := getEmployee(db, "hr")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(oneEmployee)
	fmt.Println(employees)
}

func getEmployees(db *sql.DB) ([]employee, error) {
	rows, err := db.Query("select * from emp")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var employees []employee
	for rows.Next() {
		var e employee
		err = rows.Scan(&e.ID, &e.Name, &e.deptId, &e.phone)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func getEmployee(db *sql.DB, id string) (employee, error) {
	var e employee
	row := db.QueryRow("SELECT * from emp WHERE id=?", id)

	err := row.Scan(&e.ID, &e.Name, &e.deptId, &e.phone)
	if err != nil {
		return employee{}, err
	}

	return e, nil
}
