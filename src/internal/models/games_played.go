package models

import (
	"time"
)

type GamesPlayed struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Title       string     `json:"title"`
	Genre       GameGenres `json:"genre"`
	Rating      int        `json:"rating"`
	Description string     `json:"description"`
}

func (GamesPlayed) TableName() string {
	return "games_played" // Use singular or your preferred name
}

type GamesPlayedDTO struct {
	Title       string     `json:"title"`
	Genre       GameGenres `json:"genre"`
	Rating      int        `json:"rating"`
	Description string     `json:"description"`
}

func ToGamesPlayedDTO(gamesPlayed *GamesPlayed) *GamesPlayedDTO {
	return &GamesPlayedDTO{
		Title:       gamesPlayed.Title,
		Genre:       gamesPlayed.Genre,
		Rating:      gamesPlayed.Rating,
		Description: gamesPlayed.Description,
	}
}

func ToGamesPlayedListDTO(gamesPlayed []*GamesPlayed) []*GamesPlayedDTO {
	var gamesPlayedDTOList []*GamesPlayedDTO
	for _, game := range gamesPlayed {
		dto := ToGamesPlayedDTO(game)
		gamesPlayedDTOList = append(gamesPlayedDTOList, dto)
	}
	return gamesPlayedDTOList
}
