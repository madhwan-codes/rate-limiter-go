package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"message"`
}

func EndPointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "success",
		Body:   "hello world and pong",
	}
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/ping", EndPointHandler)
	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("There was some error on listening on port 8080 : ", err)
	}
}
