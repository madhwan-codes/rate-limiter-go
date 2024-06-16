package main

import (
	"encoding/json"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "Request Failed",
				Body:   "The API's rate limit exceeded",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(message)
			if err != nil {
				panic(err)
			}
			return
		} else {
			next(w, r)
		}
	})
}
