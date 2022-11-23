package Assignement1

import (
	"encoding/json"
	"io"
	"log"
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
	var employees []Employee

	respBody, _ := json.Marshal(employee)
	_, _ = writer.Write(respBody)
	if request.Method == http.MethodGet {
		var response string
		for _, value := range employees {
			response += "ID: " + value.Id + "\tName: " + value.Name + "\n"
		}

		if response == "" {
			response = "No employees"
		}

		_, err := io.WriteString(writer, response)
		if err != nil {
			log.Println(err)
			return
		}
	}

	return
}
func handlerPost(writer http.ResponseWriter, request *http.Request) {
	var employees []Employee
	if request.Method == http.MethodPost {
		var employee Employee
		respBody, _ := json.Marshal(employee)
		_, err := writer.Write(respBody)
		if err != nil {
			return
		}
		employees = append(employees, employee)
		_, err = io.WriteString(writer, "Data added successfully")
		if err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/employee", handlerGet)     //endpoint and function
	http.HandleFunc("/employee", handlerPost)    //endpoint and function
	log.Fatal(http.ListenAndServe(":8080", nil)) //start the server
}
