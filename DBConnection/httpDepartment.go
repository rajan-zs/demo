package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var Db *sql.DB

type Department struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func startserver() {
	db, err := sql.Open("mysql", "rajan:mypassword@tcp(127.0.0.13306)/employee")
	Db = db
	if err != nil {
		return
	}
	defer Db.Close()
	err = Db.Ping()
	router := mux.NewRouter()
	router.HandleFunc("/Department", deptGetHandler).Methods("GET")
	router.HandleFunc("/Department", deptPostHandler).Methods("POST")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}

func main() {
	db, err := sql.Open("mysql", "rajan:mypassword@tcp(127.0.0.13306)/employee")
	Db = db
	if err != nil {
		return
	}
	defer Db.Close()
	err = Db.Ping()
	router := mux.NewRouter()
	router.HandleFunc("/Department", deptGetHandler).Methods("GET")
	router.HandleFunc("/Department", deptPostHandler).Methods("POST")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}

func deptPostHandler(writer http.ResponseWriter, request *http.Request) {
	var dept []Department
	var d Department
	writer.Header().Set("Content-Type", "application/json")
	req, err := ioutil.ReadAll(request.Body)
	if err != nil {
		_, err := fmt.Fprintf(writer, "enter data")
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(req, &d)
	if err != nil {
		return
	}
	dept = append(dept, d)
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(d)
	//query, err := Db.Prepare("INSERT INTO dept(id,name) VALUES(?,?)")

}

func deptGetHandler(writer http.ResponseWriter, request *http.Request) {

	var dept []Department
	rows, _ := Db.Query("SELECT * from dept")
	for rows.Next() {
		var d Department
		_ = rows.Scan(&d.Id, &d.Name)
		dept = append(dept, d)
	}
	fmt.Println(dept)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(dept)
	_, err := writer.Write(respBody)
	if err != nil {
		return
	}

}
