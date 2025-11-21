package swaggerconfig

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AddSwaggerEndpoint(router *gin.Engine) {
	// Add the Swagger route
	router.GET("/", func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")

		// If it looks like a browser request (wants HTML)
		if strings.Contains(accept, "text/html") {
			c.Redirect(http.StatusFound, "/swagger/index.html")
			return
		}

		// Otherwise, treat it as an API request
		c.JSON(http.StatusOK, gin.H{
			"status":        "API is running",
			"documentation": "/swagger/index.html",
			"version":       "1.0.0",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
