package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/khaiql/ratelimit/fixedwindow"
)

var (
	rl = fixedwindow.NewRateLimiter(3, 10*time.Second)
)

func hello(w http.ResponseWriter, req *http.Request) {
	info, err := rl.Allow("")
	if err != nil {
		log.Println(err)
	}

	if !info.Allowed {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "Rate limit exceeded. Try again in %d seconds\n", int64(info.ResetIn.Seconds()))
	} else {
		fmt.Fprintf(w, "hi!! You have %d calls left.\n", info.RemainingCalls)
	}
}

func main() {
	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8090", nil)
}
