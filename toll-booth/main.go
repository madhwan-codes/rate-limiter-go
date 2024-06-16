package main

import (
	"encoding/json"
	"fmt"
	"github.com/didip/tollbooth/v7"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
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

	message := Message{
		Status: "Request Failed",
		Body:   "The API's rate limit exceeded",
	}
	jsonMessage, _ := json.Marshal(message)

	tollBoothLimiter := tollbooth.NewLimiter(1, nil)
	tollBoothLimiter.SetMessageContentType("application/json")
	tollBoothLimiter.SetMessage(string(jsonMessage))
	http.Handle(
		"/ping",
		tollbooth.LimitFuncHandler(tollBoothLimiter, EndPointHandler),
	)
	fmt.Println("Listening on port 8000")
	http.ListenAndServe(":8080", nil)
}
