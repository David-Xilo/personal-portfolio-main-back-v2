package database

import (
	"personal-portfolio-main-back/src/internal/models"
)

type Database interface {
	GetContact() (*models.Contacts, error)
	GetProjects(projectType models.ProjectType) ([]*models.ProjectGroups, error)
	GetGamesPlayed() ([]*models.GamesPlayed, error)
}
