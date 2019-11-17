package fixedwindow

import "time"

// Storage represents an interface for storing requests made within a certain window
type Storage interface {
	CountRequest(key string, requestTs time.Time, windowDuration time.Duration) (*WindowInfo, error)
}
