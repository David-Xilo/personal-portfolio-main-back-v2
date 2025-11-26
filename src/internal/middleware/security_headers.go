package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	configuration "personal-portfolio-main-back/src/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func SecurityHeadersMiddleware(config configuration.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		isSwagger := strings.HasPrefix(path, "/swagger/") || path == "/"
		isAPIEndpoint := strings.HasPrefix(path, "/contact") || strings.HasPrefix(path, "/projects") || strings.HasPrefix(path, "/contact/") || strings.HasPrefix(path, "/projects/")
		isProd := config.Environment == "production"

		if (isProd && !isAPIEndpoint) || (!isProd && !isAPIEndpoint && !isSwagger) {
			errorMsg := fmt.Sprintf("Path not allowed %s", c.Request.URL.Path)
			slog.Error("SecurityHeadersMiddleware", "error", errorMsg)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errorMsg})

			return
		}

		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Writer.Header().Del("Server")
		c.Writer.Header().Del("X-Powered-By")
		c.Writer.Header().Del("X-AspNet-Version")
		c.Writer.Header().Del("X-AspNetMvc-Version")

		c.Header("X-Robots-Tag", "noindex, nofollow, noarchive, nosnippet, notranslate, noimageindex")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")

		// HSTS for HTTPS responses in production to enforce secure transport
		if isProd && (c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")) {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Prevent caching of sensitive responses
		if !isSwagger {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
			c.Header("Surrogate-Control", "no-store")
		}

		csp := getCSPPolicy(isSwagger)
		c.Header("Content-Security-Policy", csp)

		c.Next()
	}
}

func getCSPPolicy(isSwagger bool) string {
	if isSwagger {
		return "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data:; " +
			"font-src 'self'; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none';"
	}

	return "default-src 'none'; frame-ancestors 'none';"
}
