package rate

import (
	"sync"
	"time"
)

type RateLimiter struct {
	interval  time.Duration
	limit     int
	requests  int
	lastCheck time.Time
	mu        sync.Mutex
}

func NewRateLimiter(interval time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		interval:  interval,
		limit:     limit,
		requests:  0,
		lastCheck: time.Now(),
	}
}

func (rl *RateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if time.Since(rl.lastCheck) > rl.interval {
		rl.requests = 0
		rl.lastCheck = time.Now()
	}

	if rl.requests >= rl.limit {
		return true
	}

	rl.requests++

	return false
}
