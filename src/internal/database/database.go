package database

import (
	"personal-portfolio-main-back/src/internal/models"

	"golang.org/x/net/context"
)

type Database interface {
	GetContact(context.Context) (*models.Contacts, error)
	GetProjects(ctx context.Context, projectType models.ProjectType) ([]*models.ProjectGroups, error)
}
