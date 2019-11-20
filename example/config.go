package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"
)

const (
	memoryStorage = "memory"
	redisStorage  = "redis"
)

type config struct {
	storageType string
	max         int
	window      time.Duration
}

func (cfg *config) validate() error {
	if cfg.storageType != memoryStorage && cfg.storageType != redisStorage {
		return fmt.Errorf("unexpected storage type %s, use %s or %s only", cfg.storageType, memoryStorage, redisStorage)
	}

	if cfg.max <= 0 {
		return errors.New("max must be positive")
	}

	if cfg.window < 0 {
		return errors.New("window must be positive")
	}

	return nil
}

func parseFlag() config {
	storagePtr := flag.String("storage", memoryStorage, "storage type")
	maxRequestPtr := flag.Int("max", 100, "max requests number")
	windowPtr := flag.Int("window", 3600, "limit window in second")

	flag.Parse()

	cfg := config{
		storageType: *storagePtr,
		max:         *maxRequestPtr,
		window:      time.Duration(*windowPtr) * time.Second,
	}

	if err := cfg.validate(); err != nil {
		log.Fatal(err)
	}

	return cfg
}
