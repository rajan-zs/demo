package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		f, err := template.ParseFiles("Hare krishna")
		if err != nil {
			return
		}
		err = f.Execute(w, "hello")
		if err != nil {
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
