package ratelimit

import "time"

type RequestInfo struct {
	Timestamp *time.Time
	Calls     int
}

func (r *RequestInfo) DurationFrom(from time.Time) time.Duration {
	if r.Timestamp == nil {
		return time.Duration(0)
	}

	return from.Sub(*r.Timestamp)
}
