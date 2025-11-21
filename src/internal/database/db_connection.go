package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxRetries = 15
const dbTimeoutSeconds = 30

func InitDB(config configuration.Config) *gorm.DB {

	var err error

	dsn := config.DatabaseConfig.DbUrl

	var db *gorm.DB
	for i := 0; i < maxRetries; i++ {
		db, err = attemptConnection(dsn, i+1)
		if err == nil {
			slog.Info("Connected to the database successfully")
			break
		}
		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * 2 * time.Second
			slog.Warn("Connection failed, retrying",
				"attempt", i+1,
				"max_retries", maxRetries,
				"wait_seconds", waitTime.Seconds(),
				"error", err.Error())
			time.Sleep(waitTime)
		}
	}

	if err != nil {
		slog.Error("Failed to connect to the database", "attempts", maxRetries, "error", err)
		os.Exit(1)
	}

	if err := testConnection(db); err != nil {
		slog.Error("Database connection test failed", "error", err)
		os.Exit(1)
	}

	return db
}

func attemptConnection(dsn string, attempt int) (*gorm.DB, error) {
	slog.Info("Starting database connection attempt", "attempt", attempt, "timeout in seconds", dbTimeoutSeconds)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutSeconds*time.Second)
	defer cancel()

	slog.Info("Opening GORM connection")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("GORM Open failed", "error", err, "attempt", attempt)
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}
	slog.Info("GORM connection opened; pinging with timeout")

	slog.Info("Getting underlying SQL DB")
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get underlying sql.DB", "error", err)
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	slog.Info("Configuring connection pool")
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("Pinging database")
	if err := sqlDB.PingContext(ctx); err != nil {
		slog.Error("Database ping failed", "error", err, "attempt", attempt)
		sqlDB.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	slog.Info("Database ping successful", "attempt", attempt)

	return db, nil
}

func testConnection(db *gorm.DB) error {
	var result int
	if err := db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("test query returned unexpected result: %d", result)
	}

	slog.Info("Database connection test passed")
	return nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func ValidateDBSchema(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.Contacts{}) {
		slog.Error("Database schema is outdated. Please run the migrations first.")
		os.Exit(1)
	}
}
