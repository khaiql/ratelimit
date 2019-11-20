package ratelimit

import (
	"time"
)

type RateInfo struct {
	Allowed        bool
	LastCall       time.Time
	RemainingCalls int
	ResetIn        time.Duration
}

type RateLimiter interface {
	Allow(key string) (*RateInfo, error)
	Close() error
}
