package dberrors

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"` // Don't include HTTP status in JSON response
}

var (
	ErrDatabaseTimeout = APIError{
		Code:    "DATABASE_TIMEOUT",
		Message: "The request took too long to process. Please try again.",
		Status:  http.StatusRequestTimeout,
	}

	ErrInternalServer = APIError{
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected error occurred. Please try again later.",
		Status:  http.StatusInternalServerError,
	}

	ErrNotFound = APIError{
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found.",
		Status:  http.StatusNotFound,
	}
)

func HandleDatabaseError(c *gin.Context, err error) {

	var apiError APIError

	if errors.Is(err, context.DeadlineExceeded) {
		apiError = ErrDatabaseTimeout
		slog.Warn("Database timeout occurred", "error", apiError)
	} else if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
		apiError = ErrNotFound
		slog.Info("Resource not found", "error", apiError)
	} else {
		apiError = ErrInternalServer
		slog.Error("Database error occurred", "error", apiError)
	}

	c.JSON(apiError.Status, apiError)
}
