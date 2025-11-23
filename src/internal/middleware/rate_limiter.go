package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

type IPRateLimiter struct {
	ips        map[string]*IPLimiterEntry
	mu         *sync.RWMutex
	limit      rate.Limit
	burst      int
	cleanupCtx context.Context
	cancelFunc context.CancelFunc
}

func NewIPRateLimiter(limit rate.Limit, burst int) *IPRateLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	limiter := &IPRateLimiter{
		ips:        make(map[string]*IPLimiterEntry),
		mu:         &sync.RWMutex{},
		limit:      limit,
		burst:      burst,
		cleanupCtx: ctx,
		cancelFunc: cancel,
	}

	go limiter.startCleanup()

	return limiter
}

func (ipRateLimiter *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	ipRateLimiter.mu.Lock()
	defer ipRateLimiter.mu.Unlock()

	entry, exists := ipRateLimiter.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(ipRateLimiter.limit, ipRateLimiter.burst)
		entry = &IPLimiterEntry{
			limiter:    limiter,
			lastAccess: time.Now(),
		}
		ipRateLimiter.ips[ip] = entry
	} else {
		entry.lastAccess = time.Now()
	}

	return entry.limiter
}

func (ipRateLimiter *IPRateLimiter) cleanup(maxAge time.Duration) {
	ipRateLimiter.mu.Lock()
	defer ipRateLimiter.mu.Unlock()

	now := time.Now()
	var removedCount int

	for ip, entry := range ipRateLimiter.ips {
		if now.Sub(entry.lastAccess) > maxAge {
			delete(ipRateLimiter.ips, ip)
			removedCount++
		}
	}

	if removedCount > 0 {
		slog.Debug("Rate limiter cleanup completed",
			"removed_entries", removedCount,
			"remaining_entries", len(ipRateLimiter.ips))
	}
}

func (ipRateLimiter *IPRateLimiter) startCleanup() {
	maxAge := 1 * time.Hour
	cleanupInterval := 15 * time.Minute

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ipRateLimiter.cleanup(maxAge)
		case <-ipRateLimiter.cleanupCtx.Done():
			slog.Debug("Rate limiter cleanup routine stopped")
			return
		}
	}
}

func (ipRateLimiter *IPRateLimiter) Stop() {
	if ipRateLimiter.cancelFunc != nil {
		ipRateLimiter.cancelFunc()
	}
}

func (ipRateLimiter *IPRateLimiter) GetStats() map[string]interface{} {
	ipRateLimiter.mu.RLock()
	defer ipRateLimiter.mu.RUnlock()

	return map[string]interface{}{
		"total_ips": len(ipRateLimiter.ips),
		"limit":     float64(ipRateLimiter.limit),
		"burst":     ipRateLimiter.burst,
	}
}

func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getRemoteIP(c.Request)
		lim := limiter.GetLimiter(ip)

		if !lim.Allow() {
			slog.Warn("Request rejected by IP rate limiter",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"ip", ip)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
