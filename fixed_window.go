package ratelimit

import "time"

type FixedWindow struct {
	max      int
	duration time.Duration
	storage  Storage
}

type Option func(fw *FixedWindow)

func SetStorageOption(stgOption StorageOption) Option {
	return func(fw *FixedWindow) {
		storage, ok := storageOptionMap[stgOption]
		if ok {
			fw.storage = storage
		}
	}
}

func NewFixedWindowLimiter(max int, duration time.Duration, options ...Option) *FixedWindow {
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

func (fw *FixedWindow) Allow(key string) (*RateInfo, error) {
	firstReq, err := fw.storage.GetFirstRequestInfo(key)
	if err != nil {
		return nil, err
	}

	allowed := true
	now := time.Now()

	diff := firstReq.DurationFrom(now)

	// possiblly have exceeded the rate
	if fw.max-firstReq.Calls <= 0 {
		allowed = false
	}

	// reset the counter when the request is made outside the window
	if diff > fw.duration {
		fw.storage.ResetCounter(key)
		allowed = true // make sure user can still make request if the window is reset
	}

	req, err := fw.storage.TrackRequest(key, now)
	if err != nil {
		return nil, err
	}

	remainingCalls := fw.max - req.Calls
	if remainingCalls < 0 {
		remainingCalls = 0
	}

	ri := &RateInfo{
		Allowed:              allowed,
		LastCall:             now,
		RemainingCalls:       remainingCalls,
		CounterResetInSecond: int64((fw.duration - diff).Seconds()),
	}

	return ri, nil
}
