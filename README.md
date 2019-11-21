# ratelimit

Simple RateLimit library written in Go. There are different algorithms for rate limiting, for examples:
- Token bucket
- Leaky bucket
- Fixed window counter
- Sliding window counter
- Sliding window log

This library implements a variant of Fixed Window Counter, where the window is counted since first request is made instead of the floor of current request's timestamp.

## Algorithm

Let's say each user can have maximum 100 reqs/hour. Which means:
- Window: 1 hour
- Bucket: 100 reqs

The bucket is empty after the window's duration passed since first request is made. So it's similar to 
Fixed Window counter algorithm, but the window is counted since first request is made, not the floor of current request's timestamp. 

To store the rate limit info, this library provides two different approaches:
- In memory: store in single instance memory, handle race condition
- Redis: store in a redis and can be shared across multiple instances, this is suitable for distributed environments

### Memory storage

![memory storage flow](https://i.imgur.com/oIf8bwm.jpg "Logo Title Text 1")

### Redis storage

In a distributed system with high concurency, we can use Redis storage option. 
Redis storage implementation follows the "Write-then-Read" behavior to prevent race condition. This is done by Transactional Pipeline:
```
MULTI
HINCRBY ratelimit:user_id calls 1
HSETNX ratelimit:user_id start_ts 234329823233 # timestamp of current request
HGET ratelimit:user_id start_ts
EXEC
```

The flow is similar to the diagram above, just that in case the result of command [HSETNX](https://redis.io/commands/hsetnx) returns 1, i.e. write success, we have to then set expiration of the key:

```
PEXPIREAT ratelimit:user_id expired_ts_in_millisecond
```

where `expired_ts_in_millisecond = request_ts + window_duration`

We use `PEXPIREAT` instead of `PEXPIRE` to cater for delays in previous call to Redis.

## How to use

It's pretty simple to start using the in-memory storage
```go
// memory storage is used as default if not specified
rl := fixedwindow.NewRateLimiter(max, windowDuration) 
info, _ := rl.Allow(key) // err can be ignored if using memory storage
if info.Allowed {
   // process request
} else {
   // return 429
}
```

To use redis storage, need to initialize a [redis connection pool](https://godoc.org/github.com/gomodule/redigo/redis#NewPool) and `SetStorage` option:
```go
import (
	"github.com/gomodule/redigo/redis"
	"github.com/khaiql/ratelimit"
	"github.com/khaiql/ratelimit/fixedwindow"
)

pool := &redis.Pool{
  MaxIdle:     3,
  IdleTimeout: 240 * time.Second,
  Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
}

storage := fixedwindow.NewRedisStorage(pool)
rl := fixedwindow.NewRateLimiter(max, windowDuration, fixedwindow.SetStorage(storage))

```
