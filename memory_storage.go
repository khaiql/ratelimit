package ratelimit

import (
	"sync"
	"time"
)

type MemoryStorage struct {
	mu       sync.Mutex
	countMap map[string]int
	lastReq  map[string]time.Time
	firstReq map[string]time.Time
}

func (l *MemoryStorage) GetFirstRequestInfo(key string) (*RequestInfo, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	ts, ok := l.firstReq[key]
	if !ok {
		return &RequestInfo{}, nil
	}

	calls := l.countMap[key]
	return &RequestInfo{Timestamp: &ts, Calls: calls}, nil
}

func (l *MemoryStorage) TrackRequest(key string, requestTs time.Time) (*RequestInfo, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.countMap[key] += 1
	l.lastReq[key] = requestTs
	if _, ok := l.firstReq[key]; !ok {
		l.firstReq[key] = requestTs
	}

	info := &RequestInfo{
		Timestamp: &requestTs,
		Calls:     l.countMap[key],
	}

	return info, nil
}

func (l *MemoryStorage) ResetCounter(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.countMap, key)
	delete(l.lastReq, key)
	delete(l.firstReq, key)

	return nil
}
