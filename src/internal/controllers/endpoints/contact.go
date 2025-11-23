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

type ContactController struct {
	db     database.Database
	config configuration.Config
}

func NewContactController(db database.Database, config configuration.Config) *ContactController {
	return &ContactController{
		db:     db,
		config: config,
	}
}

func (ac *ContactController) RegisterRoutes(router gin.IRouter) {
	router.GET("/contact", ac.handleContactRequest)
}

// @Summary Get contact information
// @Description Get contact information from the database
// @Tags about
// @Accept json
// @Produce json
// @Success 200 {object} models.ContactsDTO
// @Failure 404 {object} map[string]string
// @Router /contact [get]
func (ac *ContactController) handleContactRequest(ctx *gin.Context) {
	contact, err := timeout.WithTimeout(ctx.Request.Context(), ac.config.DatabaseConfig.DbTimeout, func(dbCtx context.Context) (*models.Contacts, error) {
		return ac.db.GetContact()
	})
	if err != nil {
		dberrors.HandleDatabaseError(ctx, err)
		return
	}
	contactDTO := models.ToContactsDTO(contact)
	ctx.JSON(http.StatusOK, gin.H{"message": contactDTO})
}
