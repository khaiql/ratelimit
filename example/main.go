package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	ratelimit "github.com/khaiql/ratelimiter"
)

var (
	rl = ratelimit.NewFixedWindowLimiter(3, 10*time.Second)
)

func hello(w http.ResponseWriter, req *http.Request) {
	info, err := rl.Allow("")
	if err != nil {
		log.Println(err)
	}

	if !info.Allowed {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "Rate limit exceeded. Try again in %d seconds\n", info.CounterResetInSecond)
	} else {
		fmt.Fprintf(w, "hello\n")
	}
}

func main() {
	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8090", nil)
}
