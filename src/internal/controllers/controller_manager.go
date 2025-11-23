package controllers

import (
	"net/http"
	_ "personal-portfolio-main-back/docs"
	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/controllers/endpoints"
	swaggerconfig "personal-portfolio-main-back/src/internal/controllers/swagger"
	"personal-portfolio-main-back/src/internal/database"
	"personal-portfolio-main-back/src/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Controller interface {
	RegisterRoutes(router gin.IRouter)
}

type RouterSetup struct {
	Router      *gin.Engine
	RateLimiter *middleware.IPRateLimiter
}

func SetupRoutes(db database.Database, config configuration.Config) *RouterSetup {
	controllers := getControllers(db, config)
	routerSetup := createRouter(config)

	router := routerSetup.Router

	addHealthEndpoint(router)

	protected := router.Group("/")

	registerProtectedRoutes(protected, controllers)

	swaggerconfig.AddSwaggerEndpoint(router)

	return routerSetup
}

func createRouter(config configuration.Config) *RouterSetup {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.BasicRequestValidationMiddleware())

	router.Use(middleware.SecurityHeadersMiddleware(config))

	if config.EnableHTTPSRedirect { // Railway sets this automatically
		router.Use(middleware.HttpsRedirectMiddleware())
	}

	limiter := middleware.NewIPRateLimiter(rate.Limit(5), 30)
	router.Use(middleware.RateLimitMiddleware(limiter))

	router.Use(middleware.GetCors(config))

	return &RouterSetup{
		Router:      router,
		RateLimiter: limiter,
	}
}

func getControllers(db database.Database, config configuration.Config) []Controller {
	var controllers []Controller

	contactController := endpoints.NewContactController(db, config)
	controllers = append(controllers, contactController)

	projectsController := endpoints.NewProjectsController(db, config)
	controllers = append(controllers, projectsController)

	return controllers
}

func addHealthEndpoint(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
}

func registerProtectedRoutes(router *gin.RouterGroup, controllers []Controller) {
	for _, controller := range controllers {
		controller.RegisterRoutes(router)
	}
}
