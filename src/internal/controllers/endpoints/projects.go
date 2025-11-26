package endpoints

import (
	"net/http"
	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/database"
	"personal-portfolio-main-back/src/internal/database/errors"
	"personal-portfolio-main-back/src/internal/database/timeout"
	"personal-portfolio-main-back/src/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type ProjectsController struct {
	db     database.Database
	config configuration.Config
}

func NewProjectsController(db database.Database, config configuration.Config) *ProjectsController {
	return &ProjectsController{
		db:     db,
		config: config,
	}
}

func (tc *ProjectsController) RegisterRoutes(router gin.IRouter) {
	router.GET("/projects", tc.handleProjects)
}

// @Summary Get projects
// @Description Returns a list of projects
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} []models.ProjectGroupsDTO
// @Failure 404 {object} map[string]string
// @Router /projects [get]
func (tc *ProjectsController) handleProjects(ctx *gin.Context) {
	projects, err := timeout.WithTimeout(ctx.Request.Context(), tc.config.DatabaseConfig.DbTimeout, func(dbCtx context.Context) ([]*models.ProjectGroups, error) {
		return tc.db.GetProjects(dbCtx, models.ProjectTypeTech)
	})
	if err != nil {
		dberrors.HandleDatabaseError(ctx, err)
		return
	}
	projectsDTOList := models.ToProjectGroupsDTOList(projects)
	ctx.JSON(http.StatusOK, gin.H{"message": projectsDTOList})
}
