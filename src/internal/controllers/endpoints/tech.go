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

type TechController struct {
	db     database.Database
	config configuration.Config
}

func NewTechController(db database.Database, config configuration.Config) *TechController {
	return &TechController{
		db:     db,
		config: config,
	}
}

func (tc *TechController) RegisterRoutes(router gin.IRouter) {
	router.GET("/tech/projects", tc.handleProjects)
}

// @Summary Get projects related to tech
// @Description Returns a list of tech-related projects
// @Tags tech
// @Accept  json
// @Produce  json
// @Success 200 {array} []models.ProjectGroupsDTO
// @Failure 404 {object} map[string]string
// @Router /tech/projects [get]
func (tc *TechController) handleProjects(ctx *gin.Context) {
	projects, err := timeout.WithTimeout(ctx.Request.Context(), tc.config.DatabaseConfig.DbTimeout, func(dbCtx context.Context) ([]*models.ProjectGroups, error) {
		return tc.db.GetProjects(models.ProjectTypeTech)
	})
	if err != nil {
		dberrors.HandleDatabaseError(ctx, err)
		return
	}
	projectsDTOList := models.ToProjectGroupsDTOList(projects)
	ctx.JSON(http.StatusOK, gin.H{"message": projectsDTOList})
}
