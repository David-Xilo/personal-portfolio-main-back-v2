package middleware

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

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
			slog.Warn("Request rejected by global rate limiter",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"ip", getRemoteIP(c.Request))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Global rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
