package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGamesPlayed_TableName(t *testing.T) {
	game := GamesPlayed{}
	tableName := game.TableName()

	assert.Equal(t, "games_played", tableName)
}

func TestToGamesPlayedDTO(t *testing.T) {
	// Create a test games played record
	game := &GamesPlayed{
		ID:          1,
		CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		Title:       "The Legend of Zelda",
		Genre:       "Adventure", // Assuming GameGenres is a string type
		Rating:      5,
		Description: "Amazing adventure game",
	}

	// Convert to DTO
	dto := ToGamesPlayedDTO(game)

	// Assertions
	assert.NotNil(t, dto)
	assert.Equal(t, "The Legend of Zelda", dto.Title)
	assert.Equal(t, GameGenres("Adventure"), dto.Genre)
	assert.Equal(t, 5, dto.Rating)
	assert.Equal(t, "Amazing adventure game", dto.Description)
}

func TestToGamesPlayedDTO_WithEmptyFields(t *testing.T) {
	// Create a test games played record with minimal fields
	game := &GamesPlayed{
		ID:          2,
		CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		Title:       "Simple Game",
		Genre:       "",
		Rating:      0,
		Description: "",
	}

	// Convert to DTO
	dto := ToGamesPlayedDTO(game)

	// Assertions
	assert.NotNil(t, dto)
	assert.Equal(t, "Simple Game", dto.Title)
	assert.Equal(t, GameGenres(""), dto.Genre)
	assert.Equal(t, 0, dto.Rating)
	assert.Equal(t, "", dto.Description)
}

func TestToGamesPlayedDTO_NilGame(t *testing.T) {
	// Test with nil game - this should panic in real usage
	assert.Panics(t, func() {
		ToGamesPlayedDTO(nil)
	})
}

func TestToGamesPlayedListDTO(t *testing.T) {
	// Create test games
	games := []*GamesPlayed{
		{
			ID:          1,
			Title:       "Game 1",
			Genre:       "Action",
			Rating:      5,
			Description: "First game",
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Title:       "Game 2",
			Genre:       "RPG",
			Rating:      4,
			Description: "Second game",
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	// Convert to DTO list
	dtoList := ToGamesPlayedListDTO(games)

	// Assertions
	assert.NotNil(t, dtoList)
	assert.Len(t, dtoList, 2)

	// Check first game
	assert.Equal(t, "Game 1", dtoList[0].Title)
	assert.Equal(t, GameGenres("Action"), dtoList[0].Genre)
	assert.Equal(t, 5, dtoList[0].Rating)
	assert.Equal(t, "First game", dtoList[0].Description)

	// Check second game
	assert.Equal(t, "Game 2", dtoList[1].Title)
	assert.Equal(t, GameGenres("RPG"), dtoList[1].Genre)
	assert.Equal(t, 4, dtoList[1].Rating)
	assert.Equal(t, "Second game", dtoList[1].Description)
}

func TestToGamesPlayedListDTO_EmptyList(t *testing.T) {
	// Test with empty list
	games := []*GamesPlayed{}
	dtoList := ToGamesPlayedListDTO(games)

	// Function returns a slice, which can be nil or empty
	assert.Len(t, dtoList, 0)
}

func TestToGamesPlayedListDTO_NilList(t *testing.T) {
	// Test with nil list
	dtoList := ToGamesPlayedListDTO(nil)

	// Function returns a slice, which can be nil or empty
	assert.Len(t, dtoList, 0)
}

func TestGamesPlayedStruct(t *testing.T) {
	now := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	deletedAt := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC).Add(1 * time.Hour)

	game := GamesPlayed{
		ID:          123,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   &deletedAt,
		Title:       "Test Game",
		Genre:       "Puzzle",
		Rating:      4,
		Description: "A fun puzzle game",
	}

	assert.Equal(t, uint(123), game.ID)
	assert.Equal(t, "Test Game", game.Title)
	assert.Equal(t, GameGenres("Puzzle"), game.Genre)
	assert.Equal(t, 4, game.Rating)
	assert.Equal(t, "A fun puzzle game", game.Description)
	assert.NotNil(t, game.DeletedAt)
	assert.Equal(t, deletedAt, *game.DeletedAt)
}

func TestGamesPlayedDTOStruct(t *testing.T) {
	dto := GamesPlayedDTO{
		Title:       "DTO Game",
		Genre:       "Strategy",
		Rating:      5,
		Description: "A strategic game",
	}

	assert.Equal(t, "DTO Game", dto.Title)
	assert.Equal(t, GameGenres("Strategy"), dto.Genre)
	assert.Equal(t, 5, dto.Rating)
	assert.Equal(t, "A strategic game", dto.Description)
}
