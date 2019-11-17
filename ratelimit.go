package ratelimit

import (
	"time"
)

type RateInfo struct {
	Allowed              bool
	LastCall             time.Time
	RemainingCalls       int
	CounterResetInSecond int64
}

type RateLimiter interface {
	Allow(key string) (*RateInfo, error)
}
