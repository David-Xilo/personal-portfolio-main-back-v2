package models

import "time"

type ProjectRepositories struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ProjectGroupID uint       `json:"project_group_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	LinkToGit      string     `json:"link_to_git"`
}
type RepositoriesDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	LinkToGit   string `json:"link_to_git"`
}

func ProjectRepositoriesToDTO(projectRepository *ProjectRepositories) *RepositoriesDTO {
	return &RepositoriesDTO{
		Title:       projectRepository.Title,
		Description: projectRepository.Description,
		LinkToGit:   projectRepository.LinkToGit,
	}
}
