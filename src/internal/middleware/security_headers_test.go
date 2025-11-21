package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	configuration "personal-portfolio-main-back/src/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware_BasicHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := configuration.Config{Environment: "production"}

	router := gin.New()
	router.Use(SecurityHeadersMiddleware(config))
	router.GET("/auth/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/test", nil)
	router.ServeHTTP(w, req)

	// Test basic security headers
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "no-store, no-cache, must-revalidate, private, max-age=0", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))

	// Test default CSP
	expectedCSP := "default-src 'none'; frame-ancestors 'none';"
	assert.Equal(t, expectedCSP, w.Header().Get("Content-Security-Policy"))
}

func TestSecurityHeadersMiddleware_APIEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := configuration.Config{Environment: "production"}

	apiEndpoints := []string{
		"/auth/token",
		"/about/contact",
		"/tech/projects",
		"/games/projects",
		"/finance/projects",
		"/health",
		"/internal/stats/rate-limiter",
	}

	for _, endpoint := range apiEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			router := gin.New()
			router.Use(SecurityHeadersMiddleware(config))
			router.GET(endpoint, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", endpoint, nil)
			router.ServeHTTP(w, req)

			// API endpoints should have strict CSP
			expectedCSP := "default-src 'none'; frame-ancestors 'none';"
			assert.Equal(t, expectedCSP, w.Header().Get("Content-Security-Policy"))
		})
	}
}

func TestSecurityHeadersMiddleware_SwaggerDevelopment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := configuration.Config{Environment: "development"}

	swaggerEndpoints := []string{"/swagger/", "/"}

	for _, endpoint := range swaggerEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			router := gin.New()
			router.Use(SecurityHeadersMiddleware(config))
			router.GET(endpoint, func(c *gin.Context) {
				c.String(http.StatusOK, "swagger ui")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", endpoint, nil)
			router.ServeHTTP(w, req)

			csp := w.Header().Get("Content-Security-Policy")

			// Should allow necessary directives for Swagger UI
			assert.Contains(t, csp, "default-src 'self'")
			assert.Contains(t, csp, "script-src 'self' 'unsafe-inline'")
			assert.Contains(t, csp, "style-src 'self' 'unsafe-inline'")
			assert.Contains(t, csp, "img-src 'self' data:")
			assert.Contains(t, csp, "font-src 'self'")
			assert.Contains(t, csp, "connect-src 'self'")
			assert.Contains(t, csp, "frame-ancestors 'none'")

			// Should NOT allow unsafe-eval (security improvement)
			assert.NotContains(t, csp, "'unsafe-eval'")
		})
	}
}

func TestSecurityHeadersMiddleware_SwaggerProduction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := configuration.Config{Environment: "production"}

	swaggerEndpoints := []string{"/swagger/", "/"}

	for _, endpoint := range swaggerEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			router := gin.New()
			router.Use(SecurityHeadersMiddleware(config))
			router.GET(endpoint, func(c *gin.Context) {
				c.String(http.StatusOK, "swagger ui")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", endpoint, nil)
			router.ServeHTTP(w, req)

			// In production, swagger endpoints should be blocked
			assert.Equal(t, http.StatusForbidden, w.Code)
			assert.Contains(t, w.Body.String(), "Path not allowed")
		})
	}
}

func TestGetCSPPolicy_APIEndpoint(t *testing.T) {
	csp := getCSPPolicy(false)
	expected := "default-src 'none'; frame-ancestors 'none';"
	assert.Equal(t, expected, csp)
}

func TestGetCSPPolicy_Swagger(t *testing.T) {
	csp := getCSPPolicy(true)

	assert.Contains(t, csp, "default-src 'self'")
	assert.Contains(t, csp, "script-src 'self' 'unsafe-inline'")
	assert.Contains(t, csp, "style-src 'self' 'unsafe-inline'")
	assert.Contains(t, csp, "img-src 'self' data:")
	assert.Contains(t, csp, "frame-ancestors 'none'")
	assert.NotContains(t, csp, "'unsafe-eval'")
}

func TestGetCSPPolicy_DefaultEndpoint(t *testing.T) {
	csp := getCSPPolicy(false)
	expected := "default-src 'none'; frame-ancestors 'none';"
	assert.Equal(t, expected, csp)
}

func TestSecurityHeadersMiddleware_AllEnvironments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	environments := []string{"development", "production", "staging", "test"}

	for _, env := range environments {
		t.Run(env, func(t *testing.T) {
			config := configuration.Config{Environment: env}

			router := gin.New()
			router.Use(SecurityHeadersMiddleware(config))
			router.GET("/auth/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/auth/test", nil)
			router.ServeHTTP(w, req)

			// All environments should have basic security headers
			assert.NotEmpty(t, w.Header().Get("X-Content-Type-Options"))
			assert.NotEmpty(t, w.Header().Get("X-Frame-Options"))
			assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))

			// CSP should never be empty
			csp := w.Header().Get("Content-Security-Policy")
			assert.NotEmpty(t, csp)
			assert.Contains(t, csp, "frame-ancestors 'none'")
		})
	}
}
