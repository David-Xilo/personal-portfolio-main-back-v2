package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBasicRequestValidationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicRequestValidationMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	t.Run("Valid GET request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Valid POST request with JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"test": "value"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid HTTP method", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "Method not allowed")
	})

	t.Run("URL too long", func(t *testing.T) {
		w := httptest.NewRecorder()
		longPath := "/test?" + strings.Repeat("a", 1100)
		req, _ := http.NewRequest("GET", longPath, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.Contains(t, w.Body.String(), "Request URL too large")
	})

	t.Run("Request body too large", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = MaxRequestBodySize + 1
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.Contains(t, w.Body.String(), "Request body too large")
	})

	t.Run("Missing Content-Type for POST", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"test": "value"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		// No Content-Type header
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Content-Type header required")
	})

	t.Run("Invalid Content-Type for POST", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"test": "value"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "text/plain")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
		assert.Contains(t, w.Body.String(), "Unsupported content type")
	})
}

func TestValidateHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicRequestValidationMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	t.Run("Valid headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		req.Header.Set("Accept", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid header name", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Invalid@Header", "value")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid header name")
	})

	t.Run("Header value too long", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		longValue := strings.Repeat("a", MaxHeaderLength+1)
		req.Header.Set("User-Agent", longValue)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Header value too long")
	})

	t.Run("Too many headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)

		// Add more headers than allowed
		for i := 0; i < MaxHeaderCount+1; i++ {
			req.Header.Set(fmt.Sprintf("Header-%d", i), "value")
		}
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Too many headers")
	})

	t.Run("Suspicious patterns in headers", func(t *testing.T) {
		suspiciousValues := []string{
			"<script>alert('xss')</script>",
			"javascript:alert(1)",
			"onload=alert(1)",
			"union select * from users",
			"../../../etc/passwd",
		}

		for _, value := range suspiciousValues {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("User-Agent", value)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid header value")
		}
	})

	t.Run("Control characters in headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "test\x00value")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid characters in header")
	})
}

func TestValidateURLPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicRequestValidationMiddleware())
	router.GET("/*path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	t.Run("Valid path", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/valid/path", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Suspicious patterns in path", func(t *testing.T) {
		suspiciousPaths := []string{
			"/test/<script>alert(1)</script>",
			"/test/javascript:alert(1)",
			"/test/../../../etc/passwd",
			"/test/union%20select",
			"/test/eval(",
		}

		for _, path := range suspiciousPaths {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid URL path")
		}
	})

	t.Run("Control characters in path", func(t *testing.T) {
		w := httptest.NewRecorder()
		// URL-encode the null byte to prevent HTTP parsing issues
		req, _ := http.NewRequest("GET", "/test%00path", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid characters in URL path")
	})

	t.Run("Null bytes in path", func(t *testing.T) {
		w := httptest.NewRecorder()
		// URL-encode the null byte to prevent HTTP parsing issues
		req, _ := http.NewRequest("GET", "/test%00", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid characters in URL path")
	})
}

func TestContainsSuspiciousPattern(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"normal text", false},
		{"<script>alert(1)</script>", true},
		{"javascript:alert(1)", true},
		{"vbscript:msgbox(1)", true},
		{"onload=alert(1)", true},
		{"eval(document.cookie)", true},
		{"expression(alert(1))", true},
		{"union select * from users", true},
		{"../../../etc/passwd", true},
		{"..\\..\\..\\windows\\system32", true},
		{"\\x41\\x42\\x43", true},
		{"normal-file.txt", false},
		{"user@example.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := ContainsSuspiciousPattern(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestContainsControlCharacters(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"normal text", false},
		{"text with spaces", false},
		{"text\twith\ttabs", false},
		{"text\nwith\nnewlines", false},
		{"text\rwith\rreturns", false},
		{"text\x00with\x00nulls", true},
		{"text\x01with\x01control", true},
		{"text\x7fwith\x7fdel", true},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := ContainsControlCharacters(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsValidHeaderName(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{"User-Agent", true},
		{"Content-Type", true},
		{"X-Custom-Header", true},
		{"header_with_underscores", true},
		{"Valid123", true},
		{"", false},
		{"Invalid@Header", false},
		{"Invalid Header", false},
		{"Invalid:Header", false},
		{strings.Repeat("a", 101), false}, // Too long
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidHeaderName(tc.name)
			assert.Equal(t, tc.expected, result)
		})
	}
}
