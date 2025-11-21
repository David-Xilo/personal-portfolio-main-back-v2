package endpoints

import (
	"net/http"
	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/database"
	dberrors "personal-portfolio-main-back/src/internal/database/errors"
	"personal-portfolio-main-back/src/internal/database/timeout"
	"personal-portfolio-main-back/src/internal/models"
	"personal-portfolio-main-back/src/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type AboutController struct {
	db                    database.Database
	config                configuration.Config
	personalReviewService *service.PersonalReviewService
}

func NewAboutController(db database.Database, config configuration.Config) *AboutController {
	return &AboutController{
		db:                    db,
		config:                config,
		personalReviewService: service.NewPersonalReviewService(),
	}
}

func (ac *AboutController) RegisterRoutes(router gin.IRouter) {
	router.GET("/about/contact", ac.handleContactRequest)
	router.GET("/about/reviews/carousel", ac.handleReviewsCarouselRequest)
}

// @Summary Get contact information
// @Description Get contact information from the database
// @Tags about
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ContactsDTO
// @Failure 404 {object} map[string]string
// @Router /about/contact [get]
func (ac *AboutController) handleContactRequest(ctx *gin.Context) {
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

// @Summary Get random reviews from random people
// @Description Get random reviews from random people, for the carousel component in the about section
// @Tags about
// @Accept  json
// @Produce  json
// @Success 200 {array} models.PersonalReviewsCarouselDTO
// @Failure 404 {object} map[string]string
// @Router /about/reviews/carousel [get]
func (ac *AboutController) handleReviewsCarouselRequest(c *gin.Context) {
	reviewCarousel := ac.personalReviewService.GetAllReviews()
	c.JSON(http.StatusOK, gin.H{"message": reviewCarousel})
}
