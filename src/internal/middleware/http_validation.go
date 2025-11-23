package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

const (
	MaxURLLength       = 1000
	MaxHeaderLength    = 8192
	MaxRequestBodySize = 1024 * 1024 // 1MB
	MaxHeaderCount     = 50
)

func BasicRequestValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			logSecurityEvent(c, "INVALID_HTTP_METHOD", c.Request.Method)
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Method not allowed",
			})
			c.Abort()
			return
		}

		if len(c.Request.URL.String()) > MaxURLLength {
			logSecurityEvent(c, "URL_TOO_LONG", fmt.Sprintf("URL length: %d", len(c.Request.URL.String())))
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request URL too large",
			})
			c.Abort()
			return
		}

		if c.Request.ContentLength > MaxRequestBodySize {
			logSecurityEvent(c, "REQUEST_BODY_TOO_LARGE", fmt.Sprintf("Content-Length: %d", c.Request.ContentLength))
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			c.Abort()
			return
		}

		if !validateHeaders(c) {
			return // validateHeaders handles the response
		}

		if !validateURLPath(c) {
			return // validateURLPath handles the response
		}

		c.Next()
	}
}

func validateHeaders(c *gin.Context) bool {

	if len(c.Request.Header) > MaxHeaderCount {
		logSecurityEvent(c, "TOO_MANY_HEADERS", fmt.Sprintf("Header count: %d", len(c.Request.Header)))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Too many headers",
		})
		c.Abort()
		return false
	}

	for name, values := range c.Request.Header {
		if !isValidHeaderName(name) {
			logSecurityEvent(c, "INVALID_HEADER_NAME", name)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid header name",
			})
			c.Abort()
			return false
		}

		for _, value := range values {
			if len(value) > MaxHeaderLength {
				logSecurityEvent(c, "HEADER_VALUE_TOO_LONG", fmt.Sprintf("Header: %s, Length: %d", name, len(value)))
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Header value too long",
				})
				c.Abort()
				return false
			}

			if ContainsSuspiciousPattern(value) {
				logSecurityEvent(c, "SUSPICIOUS_HEADER_VALUE", fmt.Sprintf("Header: %s, Value: %s", name, value))
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid header value",
				})
				c.Abort()
				return false
			}

			if ContainsControlCharacters(value) {
				logSecurityEvent(c, "CONTROL_CHARS_IN_HEADER", fmt.Sprintf("Header: %s", name))
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid characters in header",
				})
				c.Abort()
				return false
			}
		}
	}

	return true
}

func validateURLPath(c *gin.Context) bool {
	path := c.Request.URL.Path

	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		logSecurityEvent(c, "INVALID_URL_ENCODING", path)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL encoding",
		})
		c.Abort()
		return false
	}

	if ContainsSuspiciousPattern(path) || ContainsSuspiciousPattern(decodedPath) {
		logSecurityEvent(c, "SUSPICIOUS_URL_PATH", path)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL path",
		})
		c.Abort()
		return false
	}

	if ContainsControlCharacters(path) || ContainsControlCharacters(decodedPath) {
		logSecurityEvent(c, "CONTROL_CHARS_IN_PATH", path)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid characters in URL path",
		})
		c.Abort()
		return false
	}

	// Check for null bytes in both original and decoded path
	if strings.Contains(path, "\x00") || strings.Contains(decodedPath, "\x00") {
		logSecurityEvent(c, "NULL_BYTE_IN_PATH", path)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid characters in URL path",
		})
		c.Abort()
		return false
	}

	return true
}

func validateContentType(c *gin.Context) bool {
	contentType := c.GetHeader("Content-Type")

	allowedContentTypes := []string{
		"application/json",
	}

	if contentType == "" {
		logSecurityEvent(c, "MISSING_CONTENT_TYPE", "")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content-Type header required for POST requests",
		})
		c.Abort()
		return false
	}

	for _, allowed := range allowedContentTypes {
		if strings.HasPrefix(strings.ToLower(contentType), allowed) {
			return true
		}
	}

	logSecurityEvent(c, "INVALID_CONTENT_TYPE", contentType)
	c.JSON(http.StatusUnsupportedMediaType, gin.H{
		"error": "Unsupported content type",
	})
	c.Abort()
	return false
}

func isValidHeaderName(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '-' && char != '_' {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 100
}

func ContainsSuspiciousPattern(input string) bool {
	suspiciousPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:",
		"onload", "onclick", "onerror", "onmouseover",
		"eval(", "expression(", "document.cookie",
		"union select", "select * from", "insert into",
		"delete from", "drop table", "../", "..\\",
		"\\x", "\\u", "%2e%2e%2f", "%2e%2e%5c",
	}
	inputLower := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(inputLower, pattern) {
			return true
		}
	}
	return false
}

func ContainsControlCharacters(input string) bool {
	for _, char := range input {
		if unicode.IsControl(char) && char != '\t' && char != '\n' && char != '\r' {
			return true
		}
	}
	return false
}

func logSecurityEvent(c *gin.Context, eventType, details string) {

	path := c.Request.URL.Path
	if ContainsSuspiciousPattern(path) {
		path = "[REDACTED]"
	}

	logDetails := details
	if ContainsSuspiciousPattern(details) || ContainsControlCharacters(details) {
		logDetails = "[REDACTED]"
	}

	slog.Warn("Security validation failed",
		"event_type", eventType,
		"details", logDetails,
		"ip", getRemoteIP(c.Request),
		"path", path,
		"method", c.Request.Method,
	)
}
