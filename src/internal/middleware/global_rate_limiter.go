package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Global rate limiting (process-wide) to cap total incoming requests regardless of source IP.
const (
	globalRateLimitPerMinute = 1000 // hardcoded for now; make configurable if needed
	globalRateLimiterBurst   = 200  // allow short bursts before smoothing to the minute rate
)

type GlobalRateLimiter struct {
	limiter *rate.Limiter
}

func NewGlobalRateLimiter() *GlobalRateLimiter {
	eventsEvery := time.Minute / time.Duration(globalRateLimitPerMinute)
	return &GlobalRateLimiter{
		limiter: rate.NewLimiter(rate.Every(eventsEvery), globalRateLimiterBurst),
	}
}

func (g *GlobalRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Global rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
