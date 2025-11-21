package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToProjectGroupsDTO(t *testing.T) {
	// Create test repositories
	repos := []ProjectRepositories{
		{
			ID:             1,
			ProjectGroupID: 1,
			Title:          "Tech Repo 1",
			Description:    "First tech repository",
			LinkToGit:      "https://github.com/user/tech1",
		},
		{
			ID:             2,
			ProjectGroupID: 1,
			Title:          "Game Repo 1",
			Description:    "First game repository",
			LinkToGit:      "https://github.com/user/game1",
		},
		{
			ID:             3,
			ProjectGroupID: 1,
			Title:          "Finance Repo 1",
			Description:    "First finance repository",
			LinkToGit:      "https://github.com/user/finance1",
		},
	}

	// Create a test project group
	projectGroup := &ProjectGroups{
		ID:                  1,
		CreatedAt:           time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:           time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		Title:               "Test Project Group",
		Description:         "A test project group",
		ProjectType:         "tech",
		ProjectRepositories: repos,
	}

	// Convert to DTO
	dto := ToProjectGroupsDTO(projectGroup)

	// Assertions
	assert.NotNil(t, dto)
	assert.Equal(t, "Test Project Group", dto.Title)
	assert.Equal(t, "A test project group", dto.Description)
	assert.Equal(t, "tech", dto.ProjectType)
	assert.NotNil(t, dto.Repositories)
	assert.Len(t, dto.Repositories, 3) // Should have all repositories from all types
}

func TestToProjectGroupsDTO_EmptyRepositories(t *testing.T) {
	// Create a test project group with no repositories
	projectGroup := &ProjectGroups{
		ID:                  2,
		CreatedAt:           time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:           time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		Title:               "Empty Project Group",
		Description:         "A project group with no repositories",
		ProjectType:         "game",
		ProjectRepositories: []ProjectRepositories{},
	}

	// Convert to DTO
	dto := ToProjectGroupsDTO(projectGroup)

	// Assertions
	assert.NotNil(t, dto)
	assert.Equal(t, "Empty Project Group", dto.Title)
	assert.Equal(t, "A project group with no repositories", dto.Description)
	assert.Equal(t, "game", dto.ProjectType)
	assert.NotNil(t, dto.Repositories)
	assert.Len(t, dto.Repositories, 0)
}

func TestToProjectGroupsDTO_NilProjectGroup(t *testing.T) {
	// Test with nil project group - this should panic in real usage
	assert.Panics(t, func() {
		ToProjectGroupsDTO(nil)
	})
}

func TestToProjectGroupsDTOList(t *testing.T) {
	// Create test project groups
	projectGroups := []*ProjectGroups{
		{
			ID:          1,
			Title:       "Project Group 1",
			Description: "First project group",
			ProjectType: "tech",
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Title:       "Project Group 2",
			Description: "Second project group",
			ProjectType: "game",
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	// Convert to DTO list
	dtoList := ToProjectGroupsDTOList(projectGroups)

	// Assertions
	assert.NotNil(t, dtoList)
	assert.Len(t, dtoList, 2)

	// Check first project group
	assert.Equal(t, "Project Group 1", dtoList[0].Title)
	assert.Equal(t, "First project group", dtoList[0].Description)
	assert.Equal(t, "tech", dtoList[0].ProjectType)

	// Check second project group
	assert.Equal(t, "Project Group 2", dtoList[1].Title)
	assert.Equal(t, "Second project group", dtoList[1].Description)
	assert.Equal(t, "game", dtoList[1].ProjectType)
}

func TestToProjectGroupsDTOList_EmptyList(t *testing.T) {
	// Test with empty list
	projectGroups := []*ProjectGroups{}
	dtoList := ToProjectGroupsDTOList(projectGroups)

	assert.NotNil(t, dtoList)
	assert.Len(t, dtoList, 0)
}

func TestToProjectGroupsDTOList_NilList(t *testing.T) {
	// Test with nil list
	dtoList := ToProjectGroupsDTOList(nil)

	assert.NotNil(t, dtoList)
	assert.Len(t, dtoList, 0)
}

func TestProjectGroupsStruct(t *testing.T) {
	now := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	deletedAt := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC).Add(1 * time.Hour)

	projectGroup := ProjectGroups{
		ID:          123,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   &deletedAt,
		Title:       "Test Project",
		Description: "Test project description",
		ProjectType: "finance",
	}

	assert.Equal(t, uint(123), projectGroup.ID)
	assert.Equal(t, "Test Project", projectGroup.Title)
	assert.Equal(t, "Test project description", projectGroup.Description)
	assert.Equal(t, "finance", projectGroup.ProjectType)
	assert.NotNil(t, projectGroup.DeletedAt)
	assert.Equal(t, deletedAt, *projectGroup.DeletedAt)
}

func TestProjectGroupsDTOStruct(t *testing.T) {
	repositories := []*RepositoriesDTO{
		{
			Title:       "Test Repo",
			Description: "A test repository",
		},
	}

	dto := ProjectGroupsDTO{
		Title:        "DTO Project",
		Description:  "DTO project description",
		ProjectType:  "tech",
		Repositories: repositories,
	}

	assert.Equal(t, "DTO Project", dto.Title)
	assert.Equal(t, "DTO project description", dto.Description)
	assert.Equal(t, "tech", dto.ProjectType)
	assert.NotNil(t, dto.Repositories)
	assert.Len(t, dto.Repositories, 1)
	assert.Equal(t, "Test Repo", dto.Repositories[0].Title)
}

func TestProjectGroupsRelationships(t *testing.T) {
	// Test that the struct has the expected relationship fields
	projectGroup := ProjectGroups{}

	// These slices start as nil or empty
	assert.Len(t, projectGroup.ProjectRepositories, 0)

	// Test that we can assign to them
	projectGroup.ProjectRepositories = []ProjectRepositories{{Title: "Test 1"}, {Title: "Test 2"}}

	assert.Len(t, projectGroup.ProjectRepositories, 2)
}
