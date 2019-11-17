package fixedwindow

import (
	"time"

	ratelimit "github.com/khaiql/ratelimiter"
)

// FixedWindow represents a rate limit window
type FixedWindow struct {
	max      int
	duration time.Duration
	storage  Storage
}

// Option to be applied to a FixedWindow
type Option func(fw *FixedWindow)

func SetStorageOption(stgOption StorageOption) Option {
	return func(fw *FixedWindow) {
		storage, ok := storageOptionMap[stgOption]
		if ok {
			fw.storage = storage
		}
	}
}

// NewRateLimiter returns a new FixedWindow limiter
func NewRateLimiter(max int, duration time.Duration, options ...Option) *FixedWindow {
	defaultStorage := storageOptionMap[Memory]
	fw := &FixedWindow{
		max:      max,
		duration: duration,
		storage:  defaultStorage,
	}

	for _, opt := range options {
		opt(fw)
	}

	return fw
}

// Allow implements the RateLimiter interface
func (fw *FixedWindow) Allow(key string) (*ratelimit.RateInfo, error) {
	allowed := true
	now := time.Now()
	windowInfo, err := fw.storage.CountRequest(key, now, fw.duration)
	if err != nil {
		return nil, err
	}

	diff := now.Sub(windowInfo.StartTimestamp)

	// possiblly have exceeded the rate
	remainingCalls := fw.max - windowInfo.Calls
	if remainingCalls < 0 {
		remainingCalls = 0
		allowed = false
	}

	ri := &ratelimit.RateInfo{
		Allowed:              allowed,
		LastCall:             now,
		RemainingCalls:       remainingCalls,
		CounterResetInSecond: int64((fw.duration - diff).Seconds()),
	}

	return ri, nil
}
