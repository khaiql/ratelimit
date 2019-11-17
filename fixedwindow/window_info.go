package fixedwindow

import "time"

type WindowInfo struct {
	StartTimestamp   time.Time
	LastReqTimestamp time.Time
	Calls            int
}
