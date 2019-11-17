package fixedwindow

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khaiql/ratelimit"
)

func TestAllow_Concurrent(t *testing.T) {
	var (
		key      = "test_key"
		max      = 3
		duration = 100 * time.Millisecond
		wg       = sync.WaitGroup{}
		storage  = NewMemoryStorage()

		allowedRequests    int32
		notAllowedRequests int32

		fw = NewRateLimiter(max, duration, SetStorage(storage))
	)

	// the test simulate 5 concurrent requests
	for i := 0; i < max+2; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			result, err := fw.Allow(key)
			if err != nil {
				t.Errorf("didn't expected error but got %v", err)
			}

			if result.Allowed {
				atomic.AddInt32(&allowedRequests, 1)
			} else {
				atomic.AddInt32(&notAllowedRequests, 1)
			}
		}()
	}
	wg.Wait()

	if atomic.LoadInt32(&allowedRequests) != 3 {
		t.Errorf("expected 3 allowed requests, but got %d", atomic.LoadInt32(&allowedRequests))
	}

	if atomic.LoadInt32(&notAllowedRequests) != 2 {
		t.Errorf("expected 2 not allowed requests, but got %d", atomic.LoadInt32(&notAllowedRequests))
	}

	time.Sleep(duration)

	// Test when the request is made after the window closes
	newResult, err := fw.Allow(key)
	if !newResult.Allowed {
		t.Error("new request after the window should be allowed")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestAllow_Sequential(t *testing.T) {
	var (
		key      = "test_key"
		max      = 3
		duration = 100 * time.Millisecond
		now      = time.Now()

		fw = NewRateLimiter(max, duration)
	)

	for i := 0; i < max+3; i++ {
		mockNow := now.Add(time.Duration(i) * 10 * time.Millisecond)
		Now = func() time.Time {
			return mockNow
		}
		result, err := fw.Allow(key)
		if err != nil {
			t.Errorf("didn't expected error but got %v", err)
		}

		expectedRemainingCalls := max - i - 1
		if expectedRemainingCalls < 0 {
			expectedRemainingCalls = 0
		}
		expectedResult := ratelimit.RateInfo{
			Allowed:        i < max,
			RemainingCalls: expectedRemainingCalls,
			LastCall:       mockNow,
			ResetIn:        duration - time.Duration(i)*10*time.Millisecond,
		}

		if expectedResult.Allowed != result.Allowed {
			t.Errorf("Allowed: expected %v, got %v", expectedResult.Allowed, result.Allowed)
		}

		if expectedResult.ResetIn != result.ResetIn {
			t.Errorf("ResetIn: expected %v, got %v", expectedResult.ResetIn, result.ResetIn)
		}

		if expectedResult.LastCall != result.LastCall {
			t.Errorf("LastCall: expected %v, got %v", expectedResult.LastCall, result.LastCall)
		}

		if expectedResult.RemainingCalls != result.RemainingCalls {
			t.Errorf("RemainingCalls: expected %v, got %v", expectedResult.RemainingCalls, result.RemainingCalls)
		}
	}
}
