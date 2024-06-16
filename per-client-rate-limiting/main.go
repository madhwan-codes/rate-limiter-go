package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func ratePerClientRateLimiter(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	type Client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &Client{limiter: rate.NewLimiter(2, 4)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			message := Message{Status: "error", Body: "rate limit exceeded"}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(message)
			return
		}
		mu.Unlock()
		next(w, r)

	})
}

func endPointHandler(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)

	message := Message{
		Status: "success",
		Body:   "hello world and pong",
	}

	err := json.NewEncoder(write).Encode(message)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/ping", ratePerClientRateLimiter(endPointHandler))

	fmt.Println("server started at port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
