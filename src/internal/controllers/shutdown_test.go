package controllers

import (
	"testing"
	"time"

	configuration "personal-portfolio-main-back/src/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRouterSetup_GracefulShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(MockDatabase)
	config := configuration.Config{
		Environment:         "test",
		EnableHTTPSRedirect: false,
		Port:                "4000",
		AllowedOrigins:      []string{"http://localhost:3000"},
		DatabaseConfig: configuration.DbConfig{
			DbTimeout: 10 * time.Second,
		},
		ReadTimeout:          10 * time.Second,
		WriteTimeout:         1 * time.Second,
		JWTSigningKey:        "JWTSigningKey",
		FrontendAuthKey:      configuration.FrontendTokenAuth,
		JWTExpirationMinutes: 30,
	}

	// Create router setup
	routerSetup := SetupRoutes(mockDB, config)

	// Verify rate limiter is running
	assert.NotNil(t, routerSetup.RateLimiter)

	// Verify cleanup context is not cancelled initially
	select {
	case <-routerSetup.RateLimiter.GetCleanupContext().Done():
		t.Fatal("Cleanup context should not be cancelled initially")
	default:
		// Expected
	}

	// Test graceful shutdown
	routerSetup.RateLimiter.Stop()

	// Verify cleanup context is cancelled after Stop()
	select {
	case <-routerSetup.RateLimiter.GetCleanupContext().Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Cleanup context should be cancelled after Stop()")
	}
}
