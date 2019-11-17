package fixedwindow

import "time"

type WindowInfo struct {
	StartTimestamp   *time.Time
	LastReqTimestamp time.Time
	Calls            int
}

func (r *WindowInfo) DurationFrom(from time.Time) time.Duration {
	if r.StartTimestamp == nil {
		return time.Duration(0)
	}

	return from.Sub(*r.StartTimestamp)
}
