package ratelimit

import "time"

type StorageOption int

const (
	Memory StorageOption = iota
	Redis
)

var (
	storageOptionMap = map[StorageOption]Storage{
		Memory: &MemoryStorage{},
	}
)

type Storage interface {
	GetFirstRequestInfo(key string) (*RequestInfo, error)
	TrackRequest(key string, requestTs time.Time) (*RequestInfo, error)
	ResetCounter(key string) error
}
