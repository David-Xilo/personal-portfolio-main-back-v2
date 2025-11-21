// @title personal-portfolio
// @version 1.0
// @description personal-portfolio documentation for backend
// @termsOfService http://yourterms.com

// @contact.name API Support
// @contact.url http://www.support.com
// @contact.email support@support.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /
package controllers

import (
	"net/http"
	_ "personal-portfolio-main-back/docs"
	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/controllers/endpoints"
	swaggerconfig "personal-portfolio-main-back/src/internal/controllers/swagger"
	"personal-portfolio-main-back/src/internal/database"
	"personal-portfolio-main-back/src/internal/middleware"
	"personal-portfolio-main-back/src/internal/security"
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

func SetupRoutes(db database.Database, config configuration.Config, jwtManager *security.JWTManager) *RouterSetup {
	controllers := getControllers(db, config, jwtManager)
	routerSetup := createRouter(config)

	router := routerSetup.Router

	addHealthEndpoint(router)

	authController := endpoints.NewAuthController(config, jwtManager)
	authController.RegisterRoutes(router)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(jwtManager))

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

func getControllers(db database.Database, config configuration.Config, jwtManager *security.JWTManager) []Controller {
	var controllers []Controller

	aboutController := endpoints.NewAboutController(db, config)
	controllers = append(controllers, aboutController)

	techController := endpoints.NewTechController(db, config)
	controllers = append(controllers, techController)

	gamesController := endpoints.NewGamesController(db, config)
	controllers = append(controllers, gamesController)

	financeController := endpoints.NewFinanceController(db, config)
	controllers = append(controllers, financeController)

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
