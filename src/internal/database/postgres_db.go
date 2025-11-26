package database

import (
	"errors"
	"personal-portfolio-main-back/src/internal/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(db *gorm.DB) Database {
	return &PostgresDB{db: db}
}

func (p *PostgresDB) GetContact(dbCtx context.Context) (*models.Contacts, error) {
	var contact models.Contacts
	if err := p.db.
		WithContext(dbCtx).
		Where("active = ?", true).
		First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &contact, nil
}

func (p *PostgresDB) GetProjects(dbCtx context.Context, projectType models.ProjectType) ([]*models.ProjectGroups, error) {
	var projectGroups []*models.ProjectGroups

	if err := p.db.
		WithContext(dbCtx).
		Where("project_type = ?", projectType).
		Preload("ProjectRepositories", func(db *gorm.DB) *gorm.DB {
			return db.Order("show_priority desc, created_at desc")
		}).
		Order("show_priority desc, created_at desc").
		Find(&projectGroups).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return projectGroups, nil
}
