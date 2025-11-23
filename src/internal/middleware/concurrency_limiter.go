package middleware

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	maxConcurrentRequests = 50                    // hardcoded cap; tune per capacity
	concurrentWaitTimeout = 50 * time.Millisecond // fail fast if saturated
)

func ConcurrencyLimiterMiddleware() gin.HandlerFunc {
	sem := make(chan struct{}, maxConcurrentRequests)

	return func(c *gin.Context) {
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
			c.Next()
		case <-time.After(concurrentWaitTimeout):
			slog.Warn("Request rejected by concurrency limiter",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"ip", getRemoteIP(c.Request))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Server is busy, try again later",
			})
			c.Abort()
			return
		}
	}
}
