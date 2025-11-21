package middleware

import (
	configuration "personal-portfolio-main-back/src/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetCors(config configuration.Config) gin.HandlerFunc {
	return cors.New(getCORSConfig(config))
}

func getCORSConfig(config configuration.Config) cors.Config {

	allowedHeaders := []string{
		"content-type",
		"referer",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"user-agent",
		"x-client-version",
		"origin",
		"accept",
		"authorization",
	}

	return cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     allowedHeaders,
		AllowCredentials: true,
	}
}
