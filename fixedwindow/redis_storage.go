package fixedwindow

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisStorage struct {
	p pool
}

type pool interface {
	Get() redis.Conn
	Close() error
}

const (
	noCallsField     = "calls"
	windowStartField = "start_timestamp"
	keyPrefix        = "ratelimit:"
)

func NewRedisStorage(p pool) *RedisStorage {
	return &RedisStorage{
		p: p,
	}
}

func (r *RedisStorage) Close() error {
	return r.p.Close()
}

func (r *RedisStorage) CountRequest(key string, requestTs time.Time, windowDuration time.Duration) (*WindowInfo, error) {
	cn := r.p.Get()
	defer cn.Close()

	redisKey := r.getRedisKey(key)
	cn.Send("MULTI")
	cn.Send("HINCRBY", redisKey, noCallsField, 1)
	cn.Send("HSETNX", redisKey, windowStartField, requestTs.UnixNano())
	cn.Send("HGET", redisKey, windowStartField)

	values, err := redis.Values(cn.Do("EXEC"))
	if err != nil {
		// in case of error happen, try to delete the key
		r.deleteKey(redisKey, cn)
		return nil, err
	}

	calls, _ := redis.Int(values[0], nil)
	setStartTimeSuccess, _ := redis.Bool(values[1], nil)
	startTs, _ := redis.Int64(values[2], nil)

	// if start time is set, it means the rate limit does not exist, hence set expiry time
	if setStartTimeSuccess {
		_, err = cn.Do("PEXPIREAT", redisKey, r.keyExpiredAtInMilliSeconds(requestTs, windowDuration))
		if err != nil {
			// in case of error happen, try to delete the key
			r.deleteKey(redisKey, cn)
			return nil, err
		}
	}

	info := &WindowInfo{
		Calls:          calls,
		StartTimestamp: time.Unix(0, startTs),
	}

	return info, nil
}

func (r *RedisStorage) getRedisKey(key string) string {
	return keyPrefix + key
}

func (r *RedisStorage) deleteKey(key string, conn redis.Conn) error {
	_, err := conn.Do("DEL", key)
	return err
}

func (r *RedisStorage) keyExpiredAtInMilliSeconds(requestTs time.Time, windowDuration time.Duration) int64 {
	return requestTs.Add(windowDuration).Unix() * 1000
}
