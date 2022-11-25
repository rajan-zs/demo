package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
)

type Employee struct {
	EmpId   string `json:"EmpId"`
	EmpName string `json:"EmpName"`
	DeptId  string `json:"DeptId"`
	Phone   int    `json:"Phone"`
	DepName string `json:"depName"`
}
type Department struct {
	DepId   string `json:"DepId"`
	DepName string `json:"depName"`
}

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("mysql", "rajan:mypassword@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Println(err)
	}

	if err = db.Ping(); err != nil {
		log.Println(err)
	}
}
func main() {

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	http.HandleFunc("/employees", empHandlerGet)
	http.HandleFunc("/employee", getEmpByIdHandler)
	http.HandleFunc("/employeePost", empHandlerPost)
	http.HandleFunc("/departments", deptHandlerGet)
	http.HandleFunc("/departmentsPost", deptHandlerPost)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func getEmpByIdHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		return
	}
	id := request.URL.Query().Get("empId")
	rows := db.QueryRow("SELECT emp.empId, emp.empName, dept.deptId, emp.phone ,dept.depName FROM emp INNER JOIN dept on emp.id = dept.deptId where empId =?", id)

	var employees []Employee
	var e Employee
	err = rows.Scan(&e.EmpId, &e.EmpName, &e.DeptId, &e.Phone, &e.DepName)
	employees = append(employees, e)
	if err = rows.Err(); err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	fmt.Println(employees)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(employees)
	_, _ = writer.Write(respBody)
}

func deptHandlerGet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		return
	}
	rows, err := db.Query("SELECT * FROM dept")
	if err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var departments []Department
	for rows.Next() {
		var d Department
		err = rows.Scan(&d.DepId, &d.DepName)
		departments = append(departments, d)
	}
	if err = rows.Err(); err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	fmt.Println(departments)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(departments)
	_, _ = writer.Write(respBody)

}

//goland:noinspection ALL
func deptHandlerPost(writer http.ResponseWriter, request *http.Request) {
	var dep Department
	//var Departments []Department
	writer.Header().Set("Content-Type", "application/json")
	req, err := ioutil.ReadAll(request.Body)
	if err != nil {
		_, err := fmt.Fprintf(writer, "enter data")
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(req, &dep)
	if err != nil {
		return
	}
	query := "INSERT INTO dept(dept.deptId, dept.depName) VALUES (?, ?)"
	_, a := db.Exec(query, dep.DepId, dep.DepName)
	fmt.Println(a)
	//Departments = append(Departments, dep)

	if err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(dep)
	if err != nil {
		return
	}
}
func empHandlerGet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		return
	}
	rows, err := db.Query("SELECT emp.empId,emp.empName,dept.deptId, emp.phone ,dept.depName FROM emp INNER JOIN dept on emp.id = dept.deptId")
	//rows, err := db.Query("select  * from emp")
	if err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var employees = make([]Employee, 0)
	for rows.Next() {
		var e Employee
		err = rows.Scan(&e.EmpId, &e.EmpName, &e.DeptId, &e.Phone, &e.DepName)
		employees = append(employees, e)
	}
	if err = rows.Err(); err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	fmt.Println(employees)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(employees)
	_, _ = writer.Write(respBody)

}

//goland:noinspection ALL
func empHandlerPost(writer http.ResponseWriter, request *http.Request) {

	var emp Employee
	var Employees []Employee
	writer.Header().Set("Content-Type", "application/json")
	req, err := ioutil.ReadAll(request.Body)
	if err != nil {
		_, err := fmt.Fprintf(writer, "enter data")
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(req, &emp)
	if err != nil {
		return
	}

	query := "INSERT INTO emp(empId, empName, id, phone) VALUES (?, ?,?,?)"
	u := uuid.New()
	_, err = db.Exec(query, u, emp.EmpName, emp.DeptId, emp.Phone)
	fmt.Println(err)
	Employees = append(Employees, emp)

	if err != nil {
		http.Error(writer, http.StatusText(500), 500)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(emp)
	if err != nil {
		return
	}
}
