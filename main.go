package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	Name string
	Age  int
}

var people []Person

func main() {
	http.HandleFunc("/people", peopleHandler)
	http.HandleFunc("/health", healthCheckHandler)

	log.Println("Server start listening on port 8081")
	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func peopleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPeople(w, r)
	case http.MethodPost:
		postPerson(w, r)
	default:
		http.Error(w, "Invalid htttp method", http.StatusMethodNotAllowed)
	}
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
	fmt.Fprint(w, "get people: '%v'", people)
}

func postPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	people = append(people, person)
	fmt.Fprint(w, "post new person: '%v'", person)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "http web-server works correctly")
}
