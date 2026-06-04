package utils

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterAllowsFirstRequest(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)
	if !rl.allow("test-ip") {
		t.Fatal("expected first request to be allowed")
	}
}

func TestRateLimiterDeniesAfterBurst(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	if !rl.allow("test-ip") {
		t.Fatal("expected request 1 to be allowed")
	}
	if !rl.allow("test-ip") {
		t.Fatal("expected request 2 to be allowed")
	}
	if rl.allow("test-ip") {
		t.Fatal("expected request 3 to be denied (burst exhausted)")
	}
}

func TestRateLimiterRefillsOverTime(t *testing.T) {
	rl := NewRateLimiter(2, 50*time.Millisecond)

	if !rl.allow("test-ip") {
		t.Fatal("expected request 1 to be allowed")
	}
	if !rl.allow("test-ip") {
		t.Fatal("expected request 2 to be allowed")
	}
	if rl.allow("test-ip") {
		t.Fatal("expected request 3 to be denied")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.allow("test-ip") {
		t.Fatal("expected request after refill to be allowed")
	}
}

func TestRateLimiterMultipleKeys(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)

	if !rl.allow("ip-a") {
		t.Fatal("expected ip-a request 1 to be allowed")
	}
	if rl.allow("ip-a") {
		t.Fatal("expected ip-a request 2 to be denied")
	}

	if !rl.allow("ip-b") {
		t.Fatal("expected ip-b request 1 to be allowed")
	}
	if rl.allow("ip-b") {
		t.Fatal("expected ip-b request 2 to be denied")
	}
}

func TestRateLimiterFullRefill(t *testing.T) {
	rl := NewRateLimiter(2, 10*time.Millisecond)

	if !rl.allow("test-ip") {
		t.Fatal("expected request 1 to be allowed")
	}
	if !rl.allow("test-ip") {
		t.Fatal("expected request 2 to be allowed")
	}

	time.Sleep(30 * time.Millisecond)

	// After 3 intervals have passed, should have full burst back
	for i := 0; i < 2; i++ {
		if !rl.allow("test-ip") {
			t.Fatalf("expected request %d after full refill to be allowed", i+1)
		}
	}
}

func TestRateLimiterDoesNotExceedBurst(t *testing.T) {
	rl := NewRateLimiter(5, 10*time.Millisecond)

	// Exhaust
	for i := 0; i < 5; i++ {
		rl.allow("test-ip")
	}

	// Wait longer than burst would refill
	time.Sleep(100 * time.Millisecond)

	// Should only have burst (5) tokens back, not more
	for i := 0; i < 5; i++ {
		if !rl.allow("test-ip") {
			t.Fatalf("expected request %d after long wait to be allowed", i+1)
		}
	}

	// 6th should fail
	if rl.allow("test-ip") {
		t.Fatal("expected 6th request to be denied (should not exceed burst)")
	}
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	var wg sync.WaitGroup
	allowed := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed <- rl.allow("same-ip")
		}()
	}

	wg.Wait()
	close(allowed)

	var count int
	for a := range allowed {
		if a {
			count++
		}
	}

	if count != 1 {
		t.Fatalf("expected exactly 1 allowed under concurrency, got %d", count)
	}
}

func TestRateLimitHandlerAllowsBelowLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := RateLimit(3, time.Minute)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx.Request.RemoteAddr = "192.168.1.1:12345"

	handler(ctx)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestRateLimitHandlerReturns429(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := RateLimit(1, time.Minute)

	// First request — allowed
	resp1 := httptest.NewRecorder()
	ctx1, _ := gin.CreateTestContext(resp1)
	ctx1.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx1.Request.RemoteAddr = "192.168.1.1:12345"
	handler(ctx1)

	// Second request — should be 429
	resp2 := httptest.NewRecorder()
	ctx2, _ := gin.CreateTestContext(resp2)
	ctx2.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx2.Request.RemoteAddr = "192.168.1.1:12345"
	handler(ctx2)

	if resp2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", resp2.Code)
	}
}
