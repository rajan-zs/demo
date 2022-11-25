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

func deptHandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	rows, err := db.Query("SELECT * FROM dept")
	if err != nil {
		http.Error(w, "Not able to connect with Db", 500)
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
		http.Error(w, "Rows are not now avilable", http.StatusInternalServerError)
		return
	}
	fmt.Println(departments)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(departments)
	_, _ = w.Write(respBody)

}

func deptHandlerPost(w http.ResponseWriter, r *http.Request) {
	var dep Department
	w.Header().Set("Content-Type", "application/json")
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, err := fmt.Fprintf(w, "enter data")
		if err != nil {
			return
		}
	}
	_ = json.Unmarshal(req, &dep)
	query := "INSERT INTO dept(dept.deptId, dept.depName) VALUES (?, ?)"
	_, a := db.Exec(query, dep.DepId, dep.DepName)
	fmt.Println(a)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dep)
	if err != nil {
		return
	}
}
func empHandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT emp.empId,emp.empName,dept.deptId, emp.phone ,dept.depName FROM emp INNER JOIN dept on emp.id = dept.deptId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var employees = make([]Employee, 0)
	for rows.Next() {
		var e Employee
		err = rows.Scan(&e.EmpId, &e.EmpName, &e.DeptId, &e.Phone, &e.DepName)
		employees = append(employees, e)
	}

	w.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(employees)
	_, _ = w.Write(respBody)

}

func empHandlerPost(w http.ResponseWriter, r *http.Request) {

	var emp Employee
	var Employees []Employee
	w.Header().Set("Content-Type", "application/json")
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, err := fmt.Fprintf(w, "enter data")
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
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(emp)
	if err != nil {
		return
	}
}
