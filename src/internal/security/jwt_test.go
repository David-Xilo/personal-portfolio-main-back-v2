package security

import (
	"testing"
	"time"

	configuration "personal-portfolio-main-back/src/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateAndValidateToken(t *testing.T) {
	config := configuration.Config{
		JWTSigningKey:        "test-secret-key-for-jwt-signing",
		JWTExpirationMinutes: 30,
	}

	jwtManager := NewJWTManager(config)

	// Test token generation
	token, err := jwtManager.GenerateToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test token validation
	claims, err := jwtManager.ValidateToken(token)
	require.NoError(t, err)
	assert.NotNil(t, claims)

	// Verify claims
	assert.Equal(t, "frontend-client", claims.Subject)
	assert.Equal(t, "frontend", claims.ClientType)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	assert.True(t, claims.IssuedAt.Time.Before(time.Now().Add(1*time.Minute)))
}

func TestJWTManager_ValidateInvalidToken(t *testing.T) {
	config := configuration.Config{
		JWTSigningKey:        "test-secret-key-for-jwt-signing",
		JWTExpirationMinutes: 30,
	}

	jwtManager := NewJWTManager(config)

	// Test with invalid token
	_, err := jwtManager.ValidateToken("invalid-token")
	assert.Error(t, err)

	// Test with empty token
	_, err = jwtManager.ValidateToken("")
	assert.Error(t, err)
}

func TestJWTManager_WrongSigningKey(t *testing.T) {
	config1 := configuration.Config{
		JWTSigningKey:        "key1",
		JWTExpirationMinutes: 30,
	}
	config2 := configuration.Config{
		JWTSigningKey:        "key2",
		JWTExpirationMinutes: 30,
	}

	jwtManager1 := NewJWTManager(config1)
	jwtManager2 := NewJWTManager(config2)

	// Generate token with first manager
	token, err := jwtManager1.GenerateToken()
	require.NoError(t, err)

	// Try to validate with second manager (different key)
	_, err = jwtManager2.ValidateToken(token)
	assert.Error(t, err)
}
