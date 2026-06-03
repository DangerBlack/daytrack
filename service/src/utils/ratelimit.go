package utils

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	last    time.Time
	tokens  int
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	burst    int
	interval time.Duration
}

func NewRateLimiter(burst int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		burst:    burst,
		interval: interval,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(10 * time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.last) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[key]
	now := time.Now()

	if !ok {
		rl.visitors[key] = &visitor{last: now, tokens: rl.burst - 1}
		return true
	}

	elapsed := now.Sub(v.last)
	refill := int(elapsed / rl.interval)
	if v.tokens+refill > rl.burst {
		v.tokens = rl.burst
	} else {
		v.tokens += refill
	}
	v.last = now

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

func RateLimit(burst int, interval time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(burst, interval)
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()
		if !rl.allow(key) {
			ctx.JSON(429, gin.H{"error": "too many requests"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
