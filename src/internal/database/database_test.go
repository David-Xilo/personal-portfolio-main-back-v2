package database

import (
	"testing"

	"personal-portfolio-main-back/src/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestDatabaseInterface tests that the Database interface is properly defined
func TestDatabaseInterface(t *testing.T) {
	// This test verifies that the Database interface has the expected methods
	// We can't instantiate an interface directly, but we can verify it exists
	// and has the correct method signatures by using a mock implementation

	var db Database
	assert.Nil(t, db) // Interface should be nil when not implemented
}

// MockDatabaseForInterface is a simple implementation to test interface compliance
type MockDatabaseForInterface struct{}

func (m *MockDatabaseForInterface) GetContact() (*models.Contacts, error) {
	return nil, nil
}

func (m *MockDatabaseForInterface) GetProjects(projectType models.ProjectType) ([]*models.ProjectGroups, error) {
	return nil, nil
}

func (m *MockDatabaseForInterface) GetGamesPlayed() ([]*models.GamesPlayed, error) {
	return nil, nil
}

func TestDatabaseInterfaceImplementation(t *testing.T) {
	// Test that our mock properly implements the Database interface
	var db Database = &MockDatabaseForInterface{}

	assert.NotNil(t, db)
	assert.Implements(t, (*Database)(nil), db)

	// Test that all methods exist and can be called
	contact, err := db.GetContact()
	assert.Nil(t, contact)
	assert.Nil(t, err)

	projects, err := db.GetProjects(models.ProjectTypeTech)
	assert.Nil(t, projects)
	assert.Nil(t, err)

	games, err := db.GetGamesPlayed()
	assert.Nil(t, games)
	assert.Nil(t, err)
}

func TestDatabaseInterfaceMethodSignatures(t *testing.T) {
	// Test that the interface methods have the correct signatures
	// by verifying parameter and return types

	mock := &MockDatabaseForInterface{}

	// Test GetContact method signature
	contact, err := mock.GetContact()
	assert.IsType(t, (*models.Contacts)(nil), contact)
	assert.IsType(t, error(nil), err)

	// Test GetProjects method signature
	projects, err := mock.GetProjects(models.ProjectTypeTech)
	assert.IsType(t, ([]*models.ProjectGroups)(nil), projects)
	assert.IsType(t, error(nil), err)

	// Test GetGamesPlayed method signature
	games, err := mock.GetGamesPlayed()
	assert.IsType(t, ([]*models.GamesPlayed)(nil), games)
	assert.IsType(t, error(nil), err)
}

func TestProjectTypeUsageInInterface(t *testing.T) {
	// Test that ProjectType enum is properly used in the interface
	validProjectTypes := []models.ProjectType{
		models.ProjectTypeUndefined,
		models.ProjectTypeTech,
		models.ProjectTypeGame,
		models.ProjectTypeFinance,
	}

	mock := &MockDatabaseForInterface{}

	// Test that all project types can be passed to GetProjects
	for _, projectType := range validProjectTypes {
		projects, err := mock.GetProjects(projectType)
		assert.Nil(t, projects)
		assert.Nil(t, err)
	}
}
