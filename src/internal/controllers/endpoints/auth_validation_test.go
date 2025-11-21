package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthController_InputValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := configuration.Config{
		FrontendAuthKey:      "valid-auth-key",
		JWTExpirationMinutes: 60,
		JWTSigningKey:        "test-secret",
	}

	jwtManager := security.NewJWTManager(config)
	authController := NewAuthController(config, jwtManager)

	router := gin.New()
	authController.RegisterRoutes(router)

	t.Run("Valid authentication request", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := AuthRequest{AuthKey: "valid-auth-key"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, 3600, response.ExpiresIn)
	})

	t.Run("Empty auth key", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := AuthRequest{AuthKey: ""}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")
	})

	t.Run("Invalid auth key - wrong value", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := AuthRequest{AuthKey: "wrong-key"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authentication key")
	})

	t.Run("Suspicious patterns in auth key", func(t *testing.T) {
		suspiciousKeys := []string{
			"<script>alert('xss')</script>",
			"javascript:alert(1)",
			"onload=alert(1)",
			"union select * from users",
			"../../../etc/passwd",
			"eval(document.cookie)",
		}

		for _, key := range suspiciousKeys {
			w := httptest.NewRecorder()
			body := AuthRequest{AuthKey: key}
			jsonBody, _ := json.Marshal(body)
			req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Suspicious patterns are detected during validation and return 400
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid authentication key format")
		}
	})

	t.Run("Non-suspicious but wrong auth key", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := AuthRequest{AuthKey: "document.write"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Non-suspicious patterns pass validation but fail auth key comparison
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authentication key")
	})

	t.Run("Control characters in auth key", func(t *testing.T) {
		controlCharKeys := []string{
			"valid\x00key",
			"valid\x01key",
			"valid\x7fkey",
			"valid\x1fkey",
		}

		for _, key := range controlCharKeys {
			w := httptest.NewRecorder()
			body := AuthRequest{AuthKey: key}
			jsonBody, _ := json.Marshal(body)
			req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid characters in authentication key")
		}
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/token", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")
	})

	t.Run("Missing JSON body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/token", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")
	})
}

func TestAuthController_ValidateAuthRequest(t *testing.T) {
	config := configuration.Config{
		FrontendAuthKey:      "test-key",
		JWTExpirationMinutes: 60,
		JWTSigningKey:        "test-secret",
	}

	jwtManager := security.NewJWTManager(config)
	authController := NewAuthController(config, jwtManager)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/test", func(c *gin.Context) {
		req := AuthRequest{AuthKey: c.Query("key")}
		if authController.validateAuthRequest(&req, c) {
			c.JSON(http.StatusOK, gin.H{"valid": true})
		}
	})

	t.Run("Valid auth key", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test?key=valid-key", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Empty auth key", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test?key=", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Authentication key cannot be empty")
	})
}
