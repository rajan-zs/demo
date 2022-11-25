package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Employee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Add  string `json:"Add"`
}

var Db *sql.DB
var Emp = []Employee{
	{"101", "rajan", "Patna"},
}
var Employees = []Employee{}

func handlerGet(writer http.ResponseWriter, request *http.Request) {
	//var employees Employee
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if request.Method == http.MethodGet {
		respBody, _ := json.Marshal(Employees)
		_, _ = writer.Write(respBody)
	}
}
func handlerPost(writer http.ResponseWriter, request *http.Request) {

	var emp Employee
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
	Employees = append(Employees, emp)
	fmt.Println(Employees, emp)
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(emp)
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/employee", handlerGet)
	http.HandleFunc("/employeePost", handlerPost)
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer Db.Close()
}
