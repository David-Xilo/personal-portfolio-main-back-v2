package models

import (
	"time"
)

type ProjectGroups struct {
	ID            uint       `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ProjectType   string     `json:"project_type"`
	LinkToProject string     `json:"link_to_project"`
	ShowPriority  int        `json:"show_priority"`
	ImageUrl      string     `json:"image_url"`

	ProjectRepositories []ProjectRepositories `json:"project_repositories,omitempty" gorm:"foreignKey:ProjectGroupID"`
}

type ProjectGroupsDTO struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ProjectType   string `json:"project_type"`
	LinkToProject string `json:"link_to_project"`
	ImageUrl      string `json:"image_url"`

	Repositories []*RepositoriesDTO `json:"repositories,omitempty"`
}

func ToProjectGroupsDTO(projectGroup *ProjectGroups) *ProjectGroupsDTO {
	repositoriesDTOList := make([]*RepositoriesDTO, 0)
	for _, projectRepo := range projectGroup.ProjectRepositories {
		dto := ProjectRepositoriesToDTO(&projectRepo)
		repositoriesDTOList = append(repositoriesDTOList, dto)
	}

	return &ProjectGroupsDTO{
		Title:         projectGroup.Title,
		Description:   projectGroup.Description,
		ProjectType:   projectGroup.ProjectType,
		LinkToProject: projectGroup.LinkToProject,
		ImageUrl:      projectGroup.ImageUrl,
		Repositories:  repositoriesDTOList,
	}
}

func ToProjectGroupsDTOList(projectGroups []*ProjectGroups) []*ProjectGroupsDTO {
	projectGroupsDTOList := make([]*ProjectGroupsDTO, 0)
	for _, projectGroup := range projectGroups {
		dto := ToProjectGroupsDTO(projectGroup)
		projectGroupsDTOList = append(projectGroupsDTOList, dto)
	}
	return projectGroupsDTOList
}
