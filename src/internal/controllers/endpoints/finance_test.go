package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestFinanceController() (*FinanceController, *MockDatabase) {
	mockDB := new(MockDatabase)
	config := configuration.Config{
		Environment:         "test",
		EnableHTTPSRedirect: false,
		Port:                "4000",
		AllowedOrigins:      []string{"http://localhost:3000"},
		DatabaseConfig: configuration.DbConfig{
			DbTimeout: 10 * time.Second,
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	controller := NewFinanceController(mockDB, config)
	return controller, mockDB
}

func TestNewFinanceController(t *testing.T) {
	controller, mockDB := setupTestFinanceController()

	assert.NotNil(t, controller)
	assert.Equal(t, mockDB, controller.db)
	assert.NotNil(t, controller.config)
}

func TestFinanceController_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, _ := setupTestFinanceController()
	router := gin.New()

	controller.RegisterRoutes(router)

	// Test that routes are registered by checking the routes
	routes := router.Routes()

	// Check that the finance projects route exists
	financeRouteFound := false

	for _, route := range routes {
		if route.Path == "/finance/projects" && route.Method == "GET" {
			financeRouteFound = true
		}
	}

	assert.True(t, financeRouteFound, "Finance projects route should be registered")
}

func TestFinanceController_HandleProjects_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestFinanceController()

	// Set up mock expectation
	expectedProjects := []*models.ProjectGroups{
		{
			ID:          1,
			Title:       "Finance Project 1",
			Description: "Test finance project",
			ProjectType: string(models.ProjectTypeFinance),
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Title:       "Finance Project 2",
			Description: "Another finance project",
			ProjectType: string(models.ProjectTypeFinance),
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	mockDB.On("GetProjects", models.ProjectTypeFinance).Return(expectedProjects, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/finance/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"].([]interface{})
	assert.Len(t, message, 2)

	// Check first project
	firstProject := message[0].(map[string]interface{})
	assert.Equal(t, "Finance Project 1", firstProject["title"])
	assert.Equal(t, "Test finance project", firstProject["description"])
	assert.Equal(t, string(models.ProjectTypeFinance), firstProject["project_type"])

	mockDB.AssertExpectations(t)
}

func TestFinanceController_HandleProjects_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestFinanceController()

	// Set up mock expectation for empty result
	mockDB.On("GetProjects", models.ProjectTypeFinance).Return([]*models.ProjectGroups{}, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/finance/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"].([]interface{})
	assert.Len(t, message, 0)

	mockDB.AssertExpectations(t)
}

func TestFinanceController_HandleProjects_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestFinanceController()

	// Set up mock expectation for not found
	mockDB.On("GetProjects", models.ProjectTypeFinance).Return(nil, gorm.ErrRecordNotFound)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/finance/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockDB.AssertExpectations(t)
}

func TestFinanceController_HandleProjects_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestFinanceController()

	// Set up mock expectation for database error
	mockDB.On("GetProjects", models.ProjectTypeFinance).Return(nil, errors.New("database connection error"))

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/finance/projects", nil)
	router.ServeHTTP(w, req)

	// Should return 500 for database errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockDB.AssertExpectations(t)
}

func TestFinanceController_HandleProjects_Timeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create controller with very short timeout
	mockDB := new(MockDatabase)
	config := configuration.Config{
		Environment:         "test",
		EnableHTTPSRedirect: false,
		Port:                "4000",
		AllowedOrigins:      []string{"http://localhost:3000"},
		DatabaseConfig: configuration.DbConfig{
			DbTimeout: 1 * time.Nanosecond, // Very short timeout
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	controller := NewFinanceController(mockDB, config)

	// Set up mock expectation
	mockDB.On("GetProjects", models.ProjectTypeFinance).Return([]*models.ProjectGroups{}, nil).Maybe()

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/finance/projects", nil)
	router.ServeHTTP(w, req)

	// Should timeout and return error
	assert.Equal(t, http.StatusRequestTimeout, w.Code)
}
