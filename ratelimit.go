package ratelimit

import (
	"time"
)

// RateInfo represents rate limit status of the current request
type RateInfo struct {
	Allowed        bool
	LastCall       time.Time
	RemainingCalls int
	ResetIn        time.Duration
}

// RateLimiter is generic interface for different rate limit strategies
type RateLimiter interface {
	Allow(key string) (*RateInfo, error)
	Close() error
}
