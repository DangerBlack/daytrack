package api

import (
	"sync"
	"time"

	"512b.it/daytrack/src/models"
	"github.com/gin-gonic/gin"
)

type visitor struct {
	last    time.Time
	tokens  int
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // max requests per interval
	burst    int           // max burst
	interval time.Duration // refill interval
}

func NewRateLimiter(rate, burst int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
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

func RateLimit(rate, burst int, interval time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(rate, burst, interval)
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()
		if !rl.allow(key) {
			ctx.JSON(429, models.NewError(models.ErrorTooManyRequests, "too many requests"))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
