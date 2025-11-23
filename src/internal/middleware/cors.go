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
	allowedOriginSet := make(map[string]struct{}, len(config.AllowedOrigins))
	for _, o := range config.AllowedOrigins {
		allowedOriginSet[o] = struct{}{}
	}

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
		// Only allow explicitly configured origins; reject unknown origins.
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				// Require Origin header to be present for CORS processing.
				return false
			}
			_, ok := allowedOriginSet[origin]
			return ok
		},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     allowedHeaders,
		AllowCredentials: true,
	}
}
