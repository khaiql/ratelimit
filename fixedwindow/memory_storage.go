package fixedwindow

import (
	"sync"
	"time"
)

type MemoryStorage struct {
	mu       sync.Mutex
	countMap map[string]int
	firstReq map[string]time.Time
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		countMap: make(map[string]int),
		firstReq: make(map[string]time.Time),
	}
}

// CountRequest counts the request being made at the timestamp, with account for the window duration
// If the request is made passed the duration, a new window is start
func (l *MemoryStorage) CountRequest(key string, requestTs time.Time, windowDuration time.Duration) (*WindowInfo, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	startTs, ok := l.firstReq[key]
	// either this is the first request, or the window should have been reset
	if !ok || requestTs.Sub(startTs) > windowDuration {
		startTs = requestTs
		l.firstReq[key] = startTs
		l.countMap[key] = 0
	}

	l.countMap[key] += 1
	info := &WindowInfo{
		StartTimestamp: startTs,
		Calls:          l.countMap[key],
	}

	return info, nil
}
