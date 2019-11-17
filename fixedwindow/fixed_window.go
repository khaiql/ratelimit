package fixedwindow

import (
	"time"

	"github.com/khaiql/ratelimit"
)

// FixedWindow represents a rate limit window
type FixedWindow struct {
	max      int
	duration time.Duration
	storage  Storage
}

// Option to be applied to a FixedWindow
type Option func(fw *FixedWindow)

// SetStorage modifies the storage to be used by FixedWindow
func SetStorage(storage Storage) Option {
	return func(fw *FixedWindow) {
		fw.storage = storage
	}
}

// NewRateLimiter returns a new FixedWindow limiter
// If storage is not set, it is defaulted to MemoryStorage
func NewRateLimiter(max int, duration time.Duration, options ...Option) *FixedWindow {
	fw := &FixedWindow{
		max:      max,
		duration: duration,
		storage:  NewMemoryStorage(),
	}

	for _, opt := range options {
		opt(fw)
	}

	return fw
}

// Now is primarily used for mocking in test
var Now = func() time.Time {
	return time.Now()
}

// Allow implements the RateLimiter interface
func (fw *FixedWindow) Allow(key string) (*ratelimit.RateInfo, error) {
	allowed := true
	now := Now()
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
		Allowed:        allowed,
		LastCall:       now,
		RemainingCalls: remainingCalls,
		ResetIn:        fw.duration - diff,
	}

	return ri, nil
}
