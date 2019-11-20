package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/khaiql/ratelimit"
	"github.com/khaiql/ratelimit/fixedwindow"
)

var (
	rl ratelimit.RateLimiter
)

func main() {
	cfg := parseFlag()
	initLimiter(cfg)

	http.HandleFunc("/hello", hello)

	go func() {
		log.Println("Starting server at :8090")
		http.ListenAndServe(":8090", nil)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	rl.Close()
}

func getRedisStorage() fixedwindow.Storage {
	redisAddr := ":6379"
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		redisAddr = addr
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}

	return fixedwindow.NewRedisStorage(pool)
}

func initLimiter(cfg config) {
	var storage fixedwindow.Storage
	if cfg.storageType == redisStorage {
		storage = getRedisStorage()
	} else {
		storage = fixedwindow.NewMemoryStorage()
	}

	rl = fixedwindow.NewRateLimiter(cfg.max, cfg.window, fixedwindow.SetStorage(storage))
}
