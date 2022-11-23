package Assignement1Mux

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Employee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func handlerGet(writer http.ResponseWriter, request *http.Request) {
	employee := []Employee{
		{"INT195", "Rajan", 21},
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	respBody, _ := json.Marshal(employee)
	_, err := writer.Write(respBody)
	if err != nil {
		return
	}
	//err := json.NewDecoder(request.Body).Decode(&employee)

	//response := ""
	//for _, value := range employees {
	//	response += "ID: " + value.Id + "\tName: " + value.Name + "\n"
	//}
	//if response == "" {
	//	response = "No employees"
	//}
	//io.WriteString(writer, response)

}
func handlerPost(writer http.ResponseWriter, request *http.Request) {

	//if request.Method == "POST" {
	//	var employee Employee
	//	err := json.NewDecoder(request.Body).Decode(&employee)
	//	if err != nil {
	//		return
	//	}
	//	employees = append(employees, employee)
	//	_, err = io.WriteString(writer, "Data added successfully")
	//	if err != nil {
	//		return
	//	}
	//}
	var employees []Employee
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
	employees = append(employees, emp)
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(emp)
	if err != nil {
		return
	}
}
func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/employee", handlerGet).Methods("GET")   //endpoint and function
	router.HandleFunc("/employee", handlerPost).Methods("POST") //endpoint and function
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	} //start the server
}

func main() {
	StartServer()
}
