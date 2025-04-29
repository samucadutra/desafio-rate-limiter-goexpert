package ratelimiter

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMiniredisStorage(t *testing.T) (*miniredis.Miniredis, *RedisStorage) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	storage := NewRedisStorage(client)
	return mr, storage
}

func TestRateLimiter_AllowRequest(t *testing.T) {
	mr, storage := setupMiniredisStorage(t)
	defer mr.Close()

	rl := NewRateLimiter(storage, 3, 10*time.Second, 5*time.Second, nil, 0, true, false)

	t.Run("allow requests under limit", func(t *testing.T) {
		key := "test-key"
		limit := 3
		window := 10 * time.Second

		for i := 0; i < limit; i++ {
			allowed := rl.AllowRequest(key, limit, window)
			assert.True(t, allowed)
		}

		// Exceed limit
		allowed := rl.AllowRequest(key, limit, window)
		assert.False(t, allowed)
	})

	t.Run("allow again after window expires", func(t *testing.T) {
		key := "test-key-2"
		limit := 3
		window := 1 * time.Second

		for i := 0; i < limit; i++ {
			allowed := rl.AllowRequest(key, limit, window)
			assert.True(t, allowed)
		}

		// Exceed limit
		allowed := rl.AllowRequest(key, limit, window)
		assert.False(t, allowed)

		// Fast-forward time
		mr.FastForward(window + 100*time.Millisecond)

		// Should be allowed again
		allowed = rl.AllowRequest(key, limit, window)
		assert.True(t, allowed)
	})
}

func TestRateLimiter_BlockRequest(t *testing.T) {
	mr, storage := setupMiniredisStorage(t)
	defer mr.Close()

	rl := NewRateLimiter(storage, 3, 10*time.Second, 5*time.Second, nil, 0, true, false)

	t.Run("block request for specified duration", func(t *testing.T) {
		key := "block-key"
		blockWindow := 5 * time.Second
		limit := 4

		rl.BlockRequest(key, limit, blockWindow)

		// Request should be blocked
		allowed := rl.AllowRequest(key, 3, 10*time.Second)
		assert.False(t, allowed)

		// Fast-forward time
		mr.FastForward(blockWindow + 10*time.Second)

		// Request should be allowed again
		allowed = rl.AllowRequest(key, 3, 10*time.Second)
		assert.True(t, allowed)
	})
}
