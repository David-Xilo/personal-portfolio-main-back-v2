package database

import (
	"errors"
	"personal-portfolio-main-back/src/internal/models"

	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(db *gorm.DB) Database {
	return &PostgresDB{db: db}
}

func (p *PostgresDB) GetContact() (*models.Contacts, error) {
	var contact models.Contacts
	if err := p.db.Where("active = ?", true).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &contact, nil
}

func (p *PostgresDB) GetProjects(projectType models.ProjectType) ([]*models.ProjectGroups, error) {
	var projectGroups []*models.ProjectGroups

	if err := p.db.
		Where("project_type = ?", projectType).
		Preload("ProjectRepositories").
		Order("created_at desc").
		Find(&projectGroups).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return projectGroups, nil
}

func (p *PostgresDB) GetGamesPlayed() ([]*models.GamesPlayed, error) {
	var gamesPlayed []*models.GamesPlayed

	if err := p.db.
		Order("created_at desc").
		Limit(5). // limit for now, just in case
		Find(&gamesPlayed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.GamesPlayed{}, nil
		}
		return nil, err
	}

	return gamesPlayed, nil
}
