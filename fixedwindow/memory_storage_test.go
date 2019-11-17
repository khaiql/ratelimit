package fixedwindow

import (
	"testing"
	"time"
)

func TestCountRequest(t *testing.T) {
	tcs := []struct {
		name            string
		noRequests      int
		interval        time.Duration
		windowDuration  time.Duration
		requestTs       time.Time
		expectedCount   int
		expectedStartTs time.Time
	}{
		{
			name:            "all requests are made within the window",
			windowDuration:  100 * time.Millisecond,
			noRequests:      2,
			interval:        50 * time.Millisecond,
			requestTs:       time.Date(2019, 11, 18, 20, 00, 00, 00, time.Local),
			expectedCount:   2,
			expectedStartTs: time.Date(2019, 11, 18, 20, 00, 00, 00, time.Local),
		},
		{
			name:            "a request is made after the window",
			windowDuration:  100 * time.Millisecond,
			noRequests:      3,
			interval:        60 * time.Millisecond,
			requestTs:       time.Date(2019, 11, 18, 20, 00, 00, 00, time.Local),
			expectedCount:   1,
			expectedStartTs: time.Date(2019, 11, 18, 20, 00, 00, 00, time.Local).Add(120 * time.Millisecond),
		},
	}

	for _, tc := range tcs {
		key := "test_key"
		t.Run(tc.name, func(t *testing.T) {
			storage := newMemoryStorage()
			var window *WindowInfo
			requestTs := tc.requestTs
			for i := 0; i < tc.noRequests; i++ {
				window, _ = storage.CountRequest(key, requestTs, tc.windowDuration)
				time.Sleep(tc.interval)
				requestTs = requestTs.Add(tc.interval)
			}

			if window.Calls != tc.expectedCount {
				t.Errorf("expected %d calls, got %d", tc.expectedCount, window.Calls)
			}

			if tc.expectedStartTs != window.StartTimestamp {
				t.Errorf("expected %v, got %v", tc.expectedStartTs, window.StartTimestamp)
			}
		})
	}
}
