package main

import (
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

var Emp = []Employee{
	{"101", "rajan", "Patna"},
}

func handlerGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if request.Method == http.MethodGet {
		respBody, _ := json.Marshal(Emp)
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
	Emp = append(Emp, emp)
	writer.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(emp)
	if err != nil {
		return
	}
	_, err = writer.Write(marshal)
	if err != nil {
		return
	}
	//err = json.NewEncoder(writer).Encode(emp)

}

func main() {
	http.HandleFunc("/employee", handlerGet)
	http.HandleFunc("/employeePost", handlerPost)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
