package fixedwindow

import "time"

type StorageOption int

const (
	Memory StorageOption = iota
	Redis
)

var (
	storageOptionMap = map[StorageOption]Storage{
		Memory: newMemoryStorage(),
	}
)

type Storage interface {
	CountRequest(key string, requestTs time.Time, windowDuration time.Duration) (*WindowInfo, error)
}
