package endpoints

import (
	"net/http"
	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/database"
	dberrors "personal-portfolio-main-back/src/internal/database/errors"
	"personal-portfolio-main-back/src/internal/database/timeout"
	"personal-portfolio-main-back/src/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type FinanceController struct {
	db     database.Database
	config configuration.Config
}

func NewFinanceController(db database.Database, config configuration.Config) *FinanceController {
	return &FinanceController{
		db:     db,
		config: config,
	}
}

func (fc *FinanceController) RegisterRoutes(router gin.IRouter) {
	router.GET("/finance/projects", fc.handleProjects)
}

// @Summary Get projects related to finance
// @Description Returns a list of finance-related projects
// @Tags finance
// @Accept  json
// @Produce  json
// @Success 200 {array} []models.ProjectGroupsDTO
// @Failure 404 {object} map[string]string
// @Router /finance/projects [get]
func (fc *FinanceController) handleProjects(ctx *gin.Context) {
	projects, err := timeout.WithTimeout(ctx.Request.Context(), fc.config.DatabaseConfig.DbTimeout, func(dbCtx context.Context) ([]*models.ProjectGroups, error) {
		return fc.db.GetProjects(models.ProjectTypeFinance)
	})
	if err != nil {
		dberrors.HandleDatabaseError(ctx, err)
		return
	}
	projectsDTOList := models.ToProjectGroupsDTOList(projects)
	ctx.JSON(http.StatusOK, gin.H{"message": projectsDTOList})
}
