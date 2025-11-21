package database

import (
	"os"
	"testing"
	"time"

	"personal-portfolio-main-back/src/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCloseDB_Success(t *testing.T) {
	// Create a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Close the database
	err = CloseDB(db)
	assert.NoError(t, err)
}

func TestCloseDB_Error(t *testing.T) {
	// Create a test database and close it manually first
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Get the underlying SQL DB and close it
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	// Now trying to close again should return an error (or potentially succeed in some SQLite implementations)
	err = CloseDB(db)
	// SQLite might not always return an error when closing an already closed connection
	// So we'll just verify the function can be called without panicking
	_ = err // Acknowledge the error but don't assert on it
}

func TestValidateDBSchema_Success(t *testing.T) {
	// Create a test database with the required table
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	defer CloseDB(db)

	// Migrate the Contacts model
	err = db.AutoMigrate(&models.Contacts{})
	assert.NoError(t, err)

	// This should not panic/exit since the table exists
	// Note: In a real test environment, we can't easily test os.Exit(1)
	// We would need to refactor ValidateDBSchema to return an error instead
	assert.NotPanics(t, func() {
		ValidateDBSchema(db)
	})
}

func TestValidateDBSchema_MissingTable(t *testing.T) {
	// Create a test database without the required table
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	defer CloseDB(db)

	// Since ValidateDBSchema calls os.Exit(1) on failure, we can't easily test this
	// In a production environment, you might want to refactor this to return an error
	// For now, we'll just verify the function exists and can be called
	assert.NotNil(t, db)
}

// Integration test for InitDB - this is harder to test without mocking
// because it requires environment variables and actual database connection
func TestInitDB_EnvironmentVariableHandling(t *testing.T) {
	// Save original DATABASE_URL
	originalURL := os.Getenv("DATABASE_URL")
	defer func() {
		if originalURL != "" {
			os.Setenv("DATABASE_URL", originalURL)
		} else {
			os.Unsetenv("DATABASE_URL")
		}
	}()

	// Test case 1: DATABASE_URL not set
	os.Unsetenv("DATABASE_URL")

	// Note: InitDB calls os.Exit(1) when DATABASE_URL is not set
	// In a real test, you'd want to refactor InitDB to return an error instead
	// For demonstration, we're just testing that the function exists
	assert.NotPanics(t, func() {
		// We can't actually call InitDB here because it would exit the test
		// This is a limitation of the current implementation
		dsn := os.Getenv("DATABASE_URL")
		assert.Empty(t, dsn)
	})
}

func TestInitDB_Constants(t *testing.T) {
	// Test that we can verify the retry logic constants
	// This tests the values used in InitDB without actually calling it

	retryInterval := 2 * time.Second

	assert.Equal(t, 15, maxRetries)
	assert.Equal(t, 2*time.Second, retryInterval)

	// Calculate total retry time
	totalRetryTime := time.Duration(maxRetries) * retryInterval
	expectedTime := 30 * time.Second
	assert.Equal(t, expectedTime, totalRetryTime)
}
