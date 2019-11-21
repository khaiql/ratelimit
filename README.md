# ratelimit

[![GoDoc](https://godoc.org/github.com/khaiql/ratelimit?status.svg)](https://godoc.org/github.com/khaiql/ratelimit)

A simple RateLimit library written in Go. There are different algorithms for [rate limiting](https://en.wikipedia.org/wiki/Rate_limiting). This library implements a variant of Fixed Window Counter, where the window is counted since the first request is made instead of the floor of current request's timestamp. 

![ts](https://i.imgur.com/Nfz01gC.png "window image")

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

It's simple to start using the in-memory storage
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

## Running the example

If you want to run the example with Redis storage option, you must have Redis up and running locally. You can set redis address via environment variable:

```bash
export REDIS_ADDR=127.0.0.1:6379 
```

To start the server:
```bash
cd example
go build .
./example -max=10 -window=10 # default using memory storage, add -storage=redis if want to use Redis storage
```

Open another shell session:
```bash
curl localhost:8090/hello -H 'X-User-ID: 1'
```

Supported parameters:
```bash
./example -h
  -max int
        max requests number (default 100)
  -storage string
        storage type (default "memory")
  -window int
        limit window in second (default 3600)
```

### Load testing the example

We can quickly do load testing the example with [hey](https://github.com/rakyll/hey)
```bash
# start server with all default options
./example
```

In another session:
```bash
# test server with total 200 requests, 50 to run concurrently
hey -n 200 -c 50 -H "X-User-ID: 1" -m GET http://localhost:8090/hello
```

## Extendability

Simply implement `RateLimiter` interface for different strategies.
