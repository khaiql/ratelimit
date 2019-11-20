package fixedwindow

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func TestMain(m *testing.M) {
	code := m.Run()
	mockConn.Clear()
	os.Exit(code)
}

var mockConn = redigomock.NewConn()

type mockPool struct{}

func (p *mockPool) Get() redis.Conn {
	return mockConn
}

func (p *mockPool) Close() error {
	return mockConn.Close()
}

func TestCountRequest_FirstRequest_Success(t *testing.T) {
	storage := NewRedisStorage(&mockPool{})
	key := "test"
	redisKey := storage.getRedisKey(key)
	now := time.Now()
	duration := 1 * time.Second

	mockConn.Command("MULTI").ExpectError(nil)
	mockConn.Command("HINCRBY", redisKey, noCallsField, 1).ExpectError(nil)
	mockConn.Command("HSETNX", redisKey, windowStartField, now.UnixNano()).ExpectError(nil)
	mockConn.Command("HGET", redisKey, windowStartField).ExpectError(nil)
	mockConn.Command("EXEC").ExpectSlice(int64(1), int64(1), now.UnixNano())
	mockConn.Command("PEXPIREAT", redisKey, now.Add(duration).Unix()*1000).Expect(int64(1))

	windowInfo, err := storage.CountRequest(key, now, duration)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	if windowInfo.Calls != 1 {
		t.Errorf("expected 1 call, got %d", windowInfo.Calls)
	}

	if windowInfo.StartTimestamp.Sub(now) != 0 {
		t.Errorf("expected %v, got %v", now, windowInfo.StartTimestamp)
	}

	if mockConn.ExpectationsWereMet() != nil {
		t.Error("all expectations were not met")
	}
}

func TestCountRequest_SubsequenceRequest_Success(t *testing.T) {
	storage := NewRedisStorage(&mockPool{})
	key := "test"
	redisKey := storage.getRedisKey(key)
	startTs := time.Now()
	requestTs := startTs.Add(100 * time.Millisecond)
	duration := 1 * time.Second

	mockConn.Command("MULTI").ExpectError(nil)
	mockConn.Command("HINCRBY", redisKey, noCallsField, 1).ExpectError(nil)
	mockConn.Command("HSETNX", redisKey, windowStartField, requestTs.UnixNano()).ExpectError(nil)
	mockConn.Command("HGET", redisKey, windowStartField).ExpectError(nil)
	mockConn.Command("EXEC").ExpectSlice(int64(1), int64(0), startTs.UnixNano())

	windowInfo, err := storage.CountRequest(key, requestTs, duration)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	if windowInfo.Calls != 1 {
		t.Errorf("expected 1 call, got %d", windowInfo.Calls)
	}

	if windowInfo.StartTimestamp.Sub(startTs) != 0 {
		t.Errorf("expected %v, got %v", startTs, windowInfo.StartTimestamp)
	}

	if mockConn.ExpectationsWereMet() != nil {
		t.Error("all expectations were not met")
	}
}

func TestCountRequest_Errored(t *testing.T) {
	storage := NewRedisStorage(&mockPool{})
	key := "test"
	redisKey := storage.getRedisKey(key)
	now := time.Now()
	duration := 1 * time.Second

	mockErr := errors.New("something went wrong")
	mockConn.Command("MULTI").ExpectError(nil)
	mockConn.Command("HINCRBY", redisKey, noCallsField, 1).ExpectError(nil)
	mockConn.Command("HSETNX", redisKey, windowStartField, now.UnixNano()).ExpectError(nil)
	mockConn.Command("HGET", redisKey, windowStartField).ExpectError(nil)
	mockConn.Command("EXEC").ExpectError(mockErr)
	mockConn.Command("DEL", redisKey).ExpectError(nil)

	windowInfo, err := storage.CountRequest(key, now, duration)
	if err != mockErr {
		t.Errorf("expected %v, got %v", mockErr, err)
		t.FailNow()
	}

	if windowInfo != nil {
		t.Errorf("expected nil")
	}

	if mockConn.ExpectationsWereMet() != nil {
		t.Error("all expectations were not met")
	}
}
