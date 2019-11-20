package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	userID, err := authenticateUser(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized user\n")
		return
	}

	info, err := rl.Allow("user:" + userID)
	if err != nil {
		log.Fatal(err)
	}

	if !info.Allowed {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "Rate limit exceeded. Try again in %d seconds\n", int64(info.ResetIn.Seconds()))
	} else {
		fmt.Fprintf(w, "hi!! You have %d calls left.\n", info.RemainingCalls)
	}
}
