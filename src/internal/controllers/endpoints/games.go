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

type GamesController struct {
	db     database.Database
	config configuration.Config
}

func NewGamesController(db database.Database, config configuration.Config) *GamesController {
	return &GamesController{
		db:     db,
		config: config,
	}
}

type GamesFilter struct {
	Genre string `json:"genre"`
}

func (gc *GamesController) RegisterRoutes(router gin.IRouter) {
	router.GET("/games/projects", gc.handleProjects)
	router.GET("/games/played/carousel", gc.handleGamesPlayedCarousel)
}

// @Summary Get projects related to games
// @Description Returns a list of projects related to games
// @Tags games
// @Accept  json
// @Produce  json
// @Success 200 {array} []models.ProjectGroupsDTO
// @Failure 404 {object} map[string]string
// @Router /games/projects [get]
func (gc *GamesController) handleProjects(ctx *gin.Context) {
	games, err := timeout.WithTimeout(ctx.Request.Context(), gc.config.DatabaseConfig.DbTimeout, func(dbCtx context.Context) ([]*models.ProjectGroups, error) {
		return gc.db.GetProjects(models.ProjectTypeGame)
	})
	if err != nil {
		dberrors.HandleDatabaseError(ctx, err)
		return
	}
	projectsDTOList := models.ToProjectGroupsDTOList(games)
	ctx.JSON(http.StatusOK, gin.H{"message": projectsDTOList})
}

// @Summary Get projects related to games
// @Description Returns a list of projects related to games
// @Tags games
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.GamesPlayedDTO
// @Failure 404 {object} map[string]string
// @Router /games/projects [get]
func (gc *GamesController) handleGamesPlayedCarousel(ctx *gin.Context) {
	gamesPlayed, err := gc.db.GetGamesPlayed()
	if err != nil {
		dberrors.HandleDatabaseError(ctx, err)
		return
	}
	firstFive := gamesPlayed[:min(len(gamesPlayed), 5)]
	gamesPlayedDTOList := models.ToGamesPlayedListDTO(firstFive)
	ctx.JSON(http.StatusOK, gin.H{"message": gamesPlayedDTOList})
}
