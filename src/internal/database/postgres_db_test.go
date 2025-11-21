package database

import (
	"testing"
	"time"

	"personal-portfolio-main-back/src/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForPostgres(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Contacts{},
		&models.ProjectGroups{},
		&models.ProjectRepositories{},
	)
	require.NoError(t, err)

	return db
}

func TestNewPostgresDB(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	assert.NotNil(t, postgresDB)
	assert.Implements(t, (*Database)(nil), postgresDB)
}

func TestPostgresDB_GetContact_Success(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	// Create test contact
	testContact := &models.Contacts{
		Name:     "John Doe",
		Email:    "john@example.com",
		Active:   true,
		Linkedin: "linkedin.com/in/johndoe",
		Github:   "github.com/johndoe",
		Credly:   "credly.com/johndoe",
	}

	result := db.Create(testContact)
	require.NoError(t, result.Error)

	contact, err := postgresDB.GetContact()

	assert.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, "John Doe", contact.Name)
	assert.Equal(t, "john@example.com", contact.Email)
}

func TestPostgresDB_GetContact_NotFound(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	// Don't create any active records - the contacts model already has an active field

	contact, err := postgresDB.GetContact()

	assert.Error(t, err)
	assert.Nil(t, contact)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPostgresDB_GetProjects_Success(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	// Create test project group
	testProjectGroup := &models.ProjectGroups{
		Title:       "Test Project",
		Description: "Test Description",
		ProjectType: string(models.ProjectTypeTech),
		CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
	}

	result := db.Create(testProjectGroup)
	require.NoError(t, result.Error)

	projects, err := postgresDB.GetProjects(models.ProjectTypeTech)

	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "Test Project", projects[0].Title)
	assert.Equal(t, string(models.ProjectTypeTech), projects[0].ProjectType)
}

func TestPostgresDB_GetProjects_EmptyResult(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	projects, err := postgresDB.GetProjects(models.ProjectTypeTech)

	assert.NoError(t, err)
	assert.Empty(t, projects)
}

func TestPostgresDB_GetProjects_FiltersByType(t *testing.T) {
	db := setupTestDBForPostgres(t)
	postgresDB := NewPostgresDB(db)

	// Create projects of different types
	techProject := &models.ProjectGroups{
		Title:       "Tech Project",
		Description: "Tech Description",
		ProjectType: string(models.ProjectTypeTech),
		CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
	}

	gameProject := &models.ProjectGroups{
		Title:       "Game Project",
		Description: "Game Description",
		ProjectType: string(models.ProjectTypeGame),
		CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
	}

	db.Create(techProject)
	db.Create(gameProject)

	// Test tech projects
	techProjects, err := postgresDB.GetProjects(models.ProjectTypeTech)
	assert.NoError(t, err)
	assert.Len(t, techProjects, 1)
	assert.Equal(t, "Tech Project", techProjects[0].Title)

	// Test game projects
	gameProjects, err := postgresDB.GetProjects(models.ProjectTypeGame)
	assert.NoError(t, err)
	assert.Len(t, gameProjects, 1)
	assert.Equal(t, "Game Project", gameProjects[0].Title)
}

// Test for error handling - database errors should return error, not panic
func TestPostgresDB_GetContact_DatabaseError(t *testing.T) {
	db := setupTestDBForPostgres(t)

	// Close the database to simulate connection error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	postgresDB := NewPostgresDB(db)

	// Should return error instead of panicking (security fix)
	contact, err := postgresDB.GetContact()
	assert.Error(t, err)
	assert.Nil(t, contact)
	assert.Contains(t, err.Error(), "database is closed")
}
