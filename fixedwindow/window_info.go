package fixedwindow

import "time"

type WindowInfo struct {
	StartTimestamp time.Time
	Calls          int
}
