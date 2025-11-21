package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestNewIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)

	assert.NotNil(t, limiter)
	assert.NotNil(t, limiter.ips)
	assert.NotNil(t, limiter.mu)
	assert.Equal(t, rate.Limit(5), limiter.limit)
	assert.Equal(t, 10, limiter.burst)
	assert.NotNil(t, limiter.cleanupCtx)
	assert.NotNil(t, limiter.cancelFunc)

	// Cleanup
	limiter.Stop()
}

func TestIPRateLimiter_GetLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)
	defer limiter.Stop()

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	// Test getting limiter for first IP
	limiter1 := limiter.GetLimiter(ip1)
	assert.NotNil(t, limiter1)

	// Test getting limiter for same IP returns same instance
	limiter1Again := limiter.GetLimiter(ip1)
	assert.Equal(t, limiter1, limiter1Again)

	// Test getting limiter for different IP returns different instance
	limiter2 := limiter.GetLimiter(ip2)
	assert.NotNil(t, limiter2)
	// Note: We can't compare limiter instances directly as they may be equal in some implementations
	// Instead, verify that both IPs are tracked separately

	// Verify both IPs are tracked
	stats := limiter.GetStats()
	assert.Equal(t, 2, stats["total_ips"])
}

func TestIPRateLimiter_LastAccessUpdate(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)
	defer limiter.Stop()

	ip := "192.168.1.1"

	// Get limiter first time
	limiter.GetLimiter(ip)

	// Check initial last access time
	limiter.mu.RLock()
	entry1 := limiter.ips[ip]
	firstAccess := entry1.lastAccess
	limiter.mu.RUnlock()

	// Wait a bit and access again
	time.Sleep(10 * time.Millisecond)
	limiter.GetLimiter(ip)

	// Check that last access time was updated
	limiter.mu.RLock()
	entry2 := limiter.ips[ip]
	secondAccess := entry2.lastAccess
	limiter.mu.RUnlock()

	assert.True(t, secondAccess.After(firstAccess))
}

func TestIPRateLimiter_Cleanup(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)
	defer limiter.Stop()

	// Add some IPs
	limiter.GetLimiter("192.168.1.1")
	limiter.GetLimiter("192.168.1.2")
	limiter.GetLimiter("192.168.1.3")

	// Verify all IPs are tracked
	stats := limiter.GetStats()
	assert.Equal(t, 3, stats["total_ips"])

	// Manually set one entry to be old
	limiter.mu.Lock()
	limiter.ips["192.168.1.2"].lastAccess = time.Now().Add(-2 * time.Hour)
	limiter.mu.Unlock()

	// Run cleanup with 1 hour max age
	limiter.cleanup(1 * time.Hour)

	// Verify old entry was removed
	stats = limiter.GetStats()
	assert.Equal(t, 2, stats["total_ips"])

	// Verify the old IP is no longer tracked
	limiter.mu.RLock()
	_, exists := limiter.ips["192.168.1.2"]
	limiter.mu.RUnlock()
	assert.False(t, exists)
}

func TestIPRateLimiter_GetStats(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(10), 20)
	defer limiter.Stop()

	// Add some IPs
	limiter.GetLimiter("192.168.1.1")
	limiter.GetLimiter("192.168.1.2")

	stats := limiter.GetStats()

	assert.Equal(t, 2, stats["total_ips"])
	assert.Equal(t, float64(10), stats["limit"])
	assert.Equal(t, 20, stats["burst"])
}

func TestIPRateLimiter_Stop(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)

	// Verify context is not cancelled initially
	select {
	case <-limiter.cleanupCtx.Done():
		t.Fatal("Context should not be cancelled initially")
	default:
		// Expected
	}

	// Stop the limiter
	limiter.Stop()

	// Verify context is cancelled
	select {
	case <-limiter.cleanupCtx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be cancelled after Stop()")
	}
}

func TestRateLimitMiddleware_AllowsRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewIPRateLimiter(rate.Limit(10), 5) // Allow 10 req/sec, burst of 5
	defer limiter.Stop()

	router := gin.New()
	router.Use(RateLimitMiddleware(limiter))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Test that requests are allowed within limit
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimitMiddleware_BlocksExcessiveRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Very restrictive limits for testing
	limiter := NewIPRateLimiter(rate.Limit(0.1), 1) // Allow 0.1 req/sec, burst of 1
	defer limiter.Stop()

	router := gin.New()
	router.Use(RateLimitMiddleware(limiter))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// First request should be allowed (burst)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request should be blocked
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

func TestRateLimitMiddleware_DifferentIPsIndependent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewIPRateLimiter(rate.Limit(0.1), 1) // Very restrictive
	defer limiter.Stop()

	router := gin.New()
	router.Use(RateLimitMiddleware(limiter))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Request from first IP should be allowed
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Request from different IP should also be allowed (independent limits)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIPRateLimiter_MemoryLeakPrevention(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(100), 10)
	defer limiter.Stop()

	// Simulate many different IPs accessing the service
	for i := 0; i < 100; i++ {
		ip := fmt.Sprintf("192.168.1.%d", i)
		limiter.GetLimiter(ip)
	}

	// Verify all IPs are tracked
	stats := limiter.GetStats()
	assert.Equal(t, 100, stats["total_ips"])

	// Manually set half of the entries to be old
	limiter.mu.Lock()
	count := 0
	for _, entry := range limiter.ips {
		if count < 50 {
			entry.lastAccess = time.Now().Add(-2 * time.Hour)
		}
		count++
	}
	limiter.mu.Unlock()

	// Run cleanup
	limiter.cleanup(1 * time.Hour)

	// Verify old entries were removed
	stats = limiter.GetStats()
	assert.Equal(t, 50, stats["total_ips"])
}
