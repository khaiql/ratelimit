package fixedwindow

import "time"

// WindowInfo represents information of a rate limit window
type WindowInfo struct {
	StartTimestamp time.Time
	Calls          int
}
