package rate

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func Test(t *testing.T) {
	t.Run("When rate limit is exceeded", func(t *testing.T) {
		client := redis.NewClient(&redis.Options{
			//Addr: os.Getenv("REDIS_URL"),
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		defer client.Close()

		ping := client.Ping()
		if err := ping.Err(); err != nil {
			t.Errorf("Ping returned error: %s", err.Error())
		}
		limiter := NewRedisLimiter(client, 2, time.Second)

		allowed, _ := limiter.Allow("127.0.0.1")
		allowed, _ = limiter.Allow("127.0.0.1")
		if !allowed {
			t.Errorf("wanted allowed, got not allowed")
		}
		time.Sleep(1 * time.Second)
		allowed, _ = limiter.Allow("127.0.0.1")
		allowed, _ = limiter.Allow("127.0.0.1")
		allowed, _ = limiter.Allow("127.0.0.1")
		if allowed {
			t.Errorf("wanted not allowed, got allowed")
		}
	})
}
