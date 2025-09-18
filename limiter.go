package rate

import (
	"time"
)

func NewLimiter(lim Limitable, limit int, window time.Duration) *Limiter {
	return &Limiter{limitable: lim, limit: limit, window: window}
}

type Limiter struct {
	limitable Limitable
	limit     int
	window    time.Duration
	disabled  bool
}

func (rl *Limiter) Disable() {
	rl.disabled = true
}

type Limitable interface {
	Increment(key string, window time.Duration) (int64, int64, error)
	Reset(key string)
}

func (l *Limiter) Allow(ip string) (bool, error) {
	if l.disabled {
		return true, nil
	}

	count, timestamp, err := l.limitable.Increment(ip, l.window)
	if err != nil {
		return false, err
	}

	now := time.Now().Unix()
	if now-int64(l.window.Seconds()) > timestamp {
		l.limitable.Reset(ip)
		return true, nil
	}

	if count > int64(l.limit) {
		return false, nil
	}

	return true, nil
}
