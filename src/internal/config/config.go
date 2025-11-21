package configuration

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const FrontendTokenAuth = "personal-portfolio-frontend"

type Config struct {
	Environment          string
	EnableHTTPSRedirect  bool
	Port                 string
	AllowedOrigins       []string
	DatabaseConfig       DbConfig
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	JWTSigningKey        string
	FrontendAuthKey      string
	JWTExpirationMinutes int
}

type DbConfig struct {
	DbUrl     string
	DbHost    string
	DbName    string
	DbUser    string
	DbPort    string
	DbTimeout time.Duration
}

func LoadConfig() Config {
	env := GetEnvOrDefault("ENV", "development")

	isProd := env == "production"

	originList := os.Getenv("ALLOWED_ORIGINS")
	if originList == "" {
		slog.Warn("Unset FRONTEND_URL value, exiting application")
		os.Exit(1)
	}
	origins := []string{}
	for _, origin := range strings.Split(originList, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	port := GetEnvOrDefault("PORT", "4000")

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		slog.Warn("Unset DATABASE_URL value, exiting application")
		os.Exit(1)
	}

	dbTimeoutStr := GetEnvOrDefault("DATABASE_TIMEOUT", "10s")
	dbTimeout, err := time.ParseDuration(dbTimeoutStr)
	if err != nil {
		slog.Warn("Invalid DATABASE_TIMEOUT value, falling back to default", "default", "10s")
		dbTimeout = 10 * time.Second
	}

	readTimeoutStr := GetEnvOrDefault("READ_TIMEOUT", "10s")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		slog.Warn("Invalid READ_TIMEOUT value, falling back to default", "default", "10s")
		readTimeout = 10 * time.Second
	}

	// I don't have writes at the moment, used to init the server
	writeTimeoutStr := GetEnvOrDefault("WRITE_TIMEOUT", "1s")
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		slog.Warn("Invalid WRITE_TIMEOUT value, falling back to default", "default", "1s")
		writeTimeout = 1 * time.Second
	}

	// TODO - Everything related with JWT is pretty useless right now - I'll come back to it later
	jwtExpirationStr := GetEnvOrDefault("JWT_EXPIRATION_MINUTES", "30")
	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		slog.Warn("Invalid JWT_EXPIRATION_MINUTES value, falling back to default", "default", "30")
		jwtExpiration = 30
	}

	jwtSigning := GetEnvOrDefault("JWT_SIGNING_KEY", "dev_jwt_signing_key")

	dbConfig := DbConfig{
		DbUrl:     dbUrl,
		DbTimeout: dbTimeout,
	}

	return Config{
		Environment:          env,
		EnableHTTPSRedirect:  isProd,
		AllowedOrigins:       origins,
		Port:                 port,
		DatabaseConfig:       dbConfig,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		JWTSigningKey:        jwtSigning,
		FrontendAuthKey:      FrontendTokenAuth,
		JWTExpirationMinutes: jwtExpiration,
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (config *Config) IsProduction() bool {
	return config.Environment == "production"
}
