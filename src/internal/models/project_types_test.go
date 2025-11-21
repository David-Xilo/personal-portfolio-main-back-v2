package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectTypeConstants(t *testing.T) {
	// Test that constants have expected values
	assert.Equal(t, ProjectType("undefined"), ProjectTypeUndefined)
	assert.Equal(t, ProjectType("tech"), ProjectTypeTech)
	assert.Equal(t, ProjectType("game"), ProjectTypeGame)
	assert.Equal(t, ProjectType("finance"), ProjectTypeFinance)
}

func TestProjectTypeString(t *testing.T) {
	// Test string conversion
	assert.Equal(t, "undefined", string(ProjectTypeUndefined))
	assert.Equal(t, "tech", string(ProjectTypeTech))
	assert.Equal(t, "game", string(ProjectTypeGame))
	assert.Equal(t, "finance", string(ProjectTypeFinance))
}

func TestProjectTypeComparison(t *testing.T) {
	// Test equality comparisons
	assert.True(t, ProjectTypeTech == ProjectTypeTech)
	assert.False(t, ProjectTypeTech == ProjectTypeGame)
	assert.False(t, ProjectTypeGame == ProjectTypeFinance)
	assert.False(t, ProjectTypeFinance == ProjectTypeUndefined)
}

func TestProjectTypeAssignment(t *testing.T) {
	// Test variable assignment
	var projectType ProjectType

	// Default zero value should be empty string, not undefined
	assert.Equal(t, ProjectType(""), projectType)

	// Assign each type
	projectType = ProjectTypeUndefined
	assert.Equal(t, ProjectTypeUndefined, projectType)

	projectType = ProjectTypeTech
	assert.Equal(t, ProjectTypeTech, projectType)

	projectType = ProjectTypeGame
	assert.Equal(t, ProjectTypeGame, projectType)

	projectType = ProjectTypeFinance
	assert.Equal(t, ProjectTypeFinance, projectType)
}

func TestProjectTypeValidation(t *testing.T) {
	// Test valid project types
	validTypes := []ProjectType{
		ProjectTypeUndefined,
		ProjectTypeTech,
		ProjectTypeGame,
		ProjectTypeFinance,
	}

	for _, validType := range validTypes {
		assert.NotEmpty(t, string(validType))
	}

	// Test that custom values work (though they might not be valid in business logic)
	customType := ProjectType("custom")
	assert.Equal(t, "custom", string(customType))
}

func TestProjectTypeInSlice(t *testing.T) {
	// Test using project types in slices
	types := []ProjectType{
		ProjectTypeTech,
		ProjectTypeGame,
		ProjectTypeFinance,
	}

	assert.Len(t, types, 3)
	assert.Contains(t, types, ProjectTypeTech)
	assert.Contains(t, types, ProjectTypeGame)
	assert.Contains(t, types, ProjectTypeFinance)
	assert.NotContains(t, types, ProjectTypeUndefined)
}
