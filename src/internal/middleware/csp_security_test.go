package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCSPSecurityImprovements demonstrates the security improvements made to CSP policies
func TestCSPSecurityImprovements(t *testing.T) {
	// Old problematic CSP (for reference)
	oldDevelopmentCSP := "default-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';"

	// New improved CSP policies
	newSwaggerCSP := getCSPPolicy(true)
	newAPICSP := getCSPPolicy(false)

	t.Run("Swagger CSP Security Improvements", func(t *testing.T) {
		// The old policy had dangerous 'unsafe-eval' - new policy should not
		assert.Contains(t, oldDevelopmentCSP, "'unsafe-eval'", "Old policy contained unsafe-eval")
		assert.NotContains(t, newSwaggerCSP, "'unsafe-eval'", "New policy should not contain unsafe-eval")

		// New policy should still allow necessary functionality for Swagger
		assert.Contains(t, newSwaggerCSP, "'unsafe-inline'", "Should allow unsafe-inline for Swagger functionality")
		assert.Contains(t, newSwaggerCSP, "img-src 'self' data:", "Should allow data: URIs for images")
		assert.Contains(t, newSwaggerCSP, "font-src 'self'", "Should allow fonts from self")
	})

	t.Run("Swagger CSP allows necessary resources", func(t *testing.T) {
		// Swagger CSP should allow necessary resources
		assert.Contains(t, newSwaggerCSP, "'unsafe-inline'", "Should allow unsafe-inline for Swagger functionality")
		assert.NotContains(t, newSwaggerCSP, "'unsafe-eval'", "Should not allow unsafe-eval")

		// Should allow necessary resources for Swagger
		assert.Contains(t, newSwaggerCSP, "script-src 'self'", "Should allow scripts from self")
		assert.Contains(t, newSwaggerCSP, "style-src 'self'", "Should allow styles from self")
	})

	t.Run("API Endpoint CSP is Most Restrictive", func(t *testing.T) {
		// API endpoints should have the most restrictive policy
		assert.Equal(t, "default-src 'none'; frame-ancestors 'none';", newAPICSP)
		assert.NotContains(t, newAPICSP, "'self'", "API endpoints shouldn't need to load any resources")
	})

	t.Run("All Policies Block Framing", func(t *testing.T) {
		// All policies should prevent framing for clickjacking protection
		assert.Contains(t, newSwaggerCSP, "frame-ancestors 'none'")
		assert.Contains(t, newAPICSP, "frame-ancestors 'none'")
	})
}

// TestCSPVulnerabilityPrevention tests that the CSP prevents common attack vectors
func TestCSPVulnerabilityPrevention(t *testing.T) {
	testCases := []struct {
		name        string
		csp         string
		description string
	}{
		{
			name:        "API Endpoint CSP",
			csp:         getCSPPolicy(false),
			description: "Prevents all resource loading for API-only endpoints",
		},
		{
			name:        "Swagger CSP",
			csp:         getCSPPolicy(true),
			description: "Allows necessary resources but blocks dangerous eval",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// All CSP policies should prevent common attacks

			// 1. Clickjacking prevention
			assert.Contains(t, tc.csp, "frame-ancestors 'none'",
				"Should prevent clickjacking attacks")

			// 2. No unsafe-eval (prevents code injection)
			assert.NotContains(t, tc.csp, "'unsafe-eval'",
				"Should prevent JavaScript eval() based attacks")

			// 3. No external script sources (prevents external JS injection)
			assert.NotContains(t, tc.csp, "script-src *",
				"Should not allow scripts from any external source")

			// 4. No data: URIs for scripts (prevents base64 encoded malicious scripts)
			assert.NotContains(t, tc.csp, "script-src 'self' data:",
				"Should not allow data: URIs for scripts")
		})
	}
}

// TestCSPDirectiveCompleteness ensures all important security directives are covered
func TestCSPDirectiveCompleteness(t *testing.T) {
	swaggerCSP := getCSPPolicy(true)
	apiCSP := getCSPPolicy(false)

	swaggerRequiredDirectives := []string{
		"default-src",
		"script-src",
		"style-src",
		"img-src",
		"font-src",
		"connect-src",
		"frame-ancestors",
	}

	for _, directive := range swaggerRequiredDirectives {
		t.Run("Swagger_"+directive, func(t *testing.T) {
			assert.Contains(t, swaggerCSP, directive,
				"Swagger CSP should include "+directive)
		})
	}

	// API CSP should only have minimal directives
	t.Run("API_minimal_directives", func(t *testing.T) {
		assert.Contains(t, apiCSP, "default-src 'none'")
		assert.Contains(t, apiCSP, "frame-ancestors 'none'")
		assert.Equal(t, "default-src 'none'; frame-ancestors 'none';", apiCSP)
	})
}
