package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToContactsDTO(t *testing.T) {
	// Create a test contact
	contact := &Contacts{
		ID:        1,
		CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Linkedin:  "linkedin.com/in/johndoe",
		Github:    "github.com/johndoe",
		Credly:    "credly.com/johndoe",
	}

	// Convert to DTO
	dto := ToContactsDTO(contact)

	// Assertions
	assert.NotNil(t, dto)
	assert.Equal(t, "John Doe", dto.Name)
	assert.Equal(t, "john.doe@example.com", dto.Email)
	assert.Equal(t, "linkedin.com/in/johndoe", dto.LinkedIn)
	assert.Equal(t, "github.com/johndoe", dto.Github)
	assert.Equal(t, "credly.com/johndoe", dto.Credly)
}

func TestToContactsDTO_WithEmptyFields(t *testing.T) {
	// Create a test contact with empty fields
	contact := &Contacts{
		ID:        2,
		CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		Name:      "Jane Smith",
		Email:     "jane@example.com",
		Linkedin:  "",
		Github:    "",
		Credly:    "",
	}

	// Convert to DTO
	dto := ToContactsDTO(contact)

	// Assertions
	assert.NotNil(t, dto)
	assert.Equal(t, "Jane Smith", dto.Name)
	assert.Equal(t, "jane@example.com", dto.Email)
	assert.Equal(t, "", dto.LinkedIn)
	assert.Equal(t, "", dto.Github)
	assert.Equal(t, "", dto.Credly)
}

func TestToContactsDTO_NilContact(t *testing.T) {
	// Test with nil contact - this should panic in real usage
	assert.Panics(t, func() {
		ToContactsDTO(nil)
	})
}

func TestContactsStruct(t *testing.T) {
	now := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	deletedAt := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC).Add(1 * time.Hour)

	contact := Contacts{
		ID:        123,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: &deletedAt,
		Name:      "Test User",
		Email:     "test@example.com",
		Linkedin:  "linkedin.com/in/testuser",
		Github:    "github.com/testuser",
		Credly:    "credly.com/testuser",
	}

	assert.Equal(t, uint(123), contact.ID)
	assert.Equal(t, "Test User", contact.Name)
	assert.Equal(t, "test@example.com", contact.Email)
	assert.Equal(t, "linkedin.com/in/testuser", contact.Linkedin)
	assert.Equal(t, "github.com/testuser", contact.Github)
	assert.Equal(t, "credly.com/testuser", contact.Credly)
	assert.NotNil(t, contact.DeletedAt)
	assert.Equal(t, deletedAt, *contact.DeletedAt)
}

func TestContactsDTOStruct(t *testing.T) {
	dto := ContactsDTO{
		Name:     "DTO User",
		Email:    "dto@example.com",
		LinkedIn: "linkedin.com/in/dtouser",
		Github:   "github.com/dtouser",
		Credly:   "credly.com/dtouser",
	}

	assert.Equal(t, "DTO User", dto.Name)
	assert.Equal(t, "dto@example.com", dto.Email)
	assert.Equal(t, "linkedin.com/in/dtouser", dto.LinkedIn)
	assert.Equal(t, "github.com/dtouser", dto.Github)
	assert.Equal(t, "credly.com/dtouser", dto.Credly)
}
