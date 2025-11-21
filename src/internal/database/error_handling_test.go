package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Test to ensure all database methods return errors instead of panicking
func TestNoPanicOnDatabaseErrors(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(Database) error
	}{
		{
			name: "GetContact should not panic on database error",
			testFunc: func(db Database) error {
				_, err := db.GetContact()
				return err
			},
		},
		{
			name: "GetProjects should not panic on database error",
			testFunc: func(db Database) error {
				_, err := db.GetProjects("tech")
				return err
			},
		},
		{
			name: "GetGamesPlayed should not panic on database error",
			testFunc: func(db Database) error {
				_, err := db.GetGamesPlayed()
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a database connection that will fail
			db := setupTestDBForPostgres(t)

			// Close the database to simulate connection error
			sqlDB, _ := db.DB()
			sqlDB.Close()

			postgresDB := NewPostgresDB(db)

			// This should NOT panic - it should return an error
			assert.NotPanics(t, func() {
				err := tt.testFunc(postgresDB)
				assert.Error(t, err)
				assert.NotEqual(t, gorm.ErrRecordNotFound, err, "Should be a connection error, not 'not found'")
			})
		})
	}
}

// Test error handling for different types of database errors
func TestDatabaseErrorTypes(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	// Test "not found" error handling
	t.Run("GetContact handles record not found gracefully", func(t *testing.T) {
		contact, err := postgresDB.GetContact()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			assert.Nil(t, contact)
			assert.Error(t, err)
		} else {
			// If data exists, that's also fine
			assert.NoError(t, err)
			assert.NotNil(t, contact)
		}
	})

	// Test that Find operations handle empty results properly
	t.Run("GetGamesPlayed handles empty results", func(t *testing.T) {
		games, err := postgresDB.GetGamesPlayed()
		assert.NoError(t, err)  // Empty results should not error
		assert.NotNil(t, games) // Should return empty slice, not nil
	})
}
