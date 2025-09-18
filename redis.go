package rate

import (
	"time"

	"github.com/go-redis/redis"
)

func NewRedisLimiter(client *redis.Client, limit int, window time.Duration) *Limiter {
	return &Limiter{
		limitable: &RedisLimitable{client: client},
		limit:     limit,
		window:    window,
	}
}

type RedisLimitable struct {
	client *redis.Client
}

func (rl *RedisLimitable) Reset(key string) {
	rl.client.HSet(key, "timestamp", time.Now().Unix())
	rl.client.HSet(key, "count", 1)
}

func (rl *RedisLimitable) Increment(key string, window time.Duration) (int64, int64, error) {
	var (
		count     *redis.IntCmd
		timeStamp *redis.StringCmd
	)

	now := time.Now().Unix()

	_, err := rl.client.TxPipelined(func(pipe redis.Pipeliner) error {
		count = pipe.HIncrBy(key, "count", 1)
		pipe.HSetNX(key, "timestamp", now)
		pipe.Expire(key, window)
		timeStamp = pipe.HGet(key, "timestamp")
		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	c, _ := count.Result()
	ts, _ := timeStamp.Int64()

	return c, ts, nil
}
