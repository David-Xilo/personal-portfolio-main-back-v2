package middleware

import (
	"net/http"
	"time"

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
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Server is busy, try again later",
			})
			c.Abort()
			return
		}
	}
}
