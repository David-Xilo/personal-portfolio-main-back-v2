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

//// MockDatabase implements the Database interface for testing
//type MockDatabase struct {
//	mock.Mock
//}
//
//func (m *MockDatabase) GetContact() (*models.Contacts, error) {
//	args := m.Called()
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).(*models.Contacts), args.Error(1)
//}
//
//func (m *MockDatabase) GetProjects(projectType models.ProjectType) ([]*models.ProjectGroups, error) {
//	args := m.Called(projectType)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).([]*models.ProjectGroups), args.Error(1)
//}

func setupTestTechController() (*ProjectsController, *MockDatabase) {
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

	controller := NewProjectsController(mockDB, config)
	return controller, mockDB
}

func TestNewTechController(t *testing.T) {
	controller, mockDB := setupTestTechController()

	assert.NotNil(t, controller)
	assert.Equal(t, mockDB, controller.db)
	assert.NotNil(t, controller.config)
}

func TestTechController_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, _ := setupTestTechController()
	router := gin.New()

	controller.RegisterRoutes(router)

	// Test that routes are registered by checking the routes
	routes := router.Routes()

	// Check that the tech projects route exists
	techRouteFound := false

	for _, route := range routes {
		if route.Path == "/tech/projects" && route.Method == "GET" {
			techRouteFound = true
		}
	}

	assert.True(t, techRouteFound, "Tech projects route should be registered")
}

func TestTechController_HandleProjects_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestTechController()

	// Set up mock expectation
	expectedProjects := []*models.ProjectGroups{
		{
			ID:          1,
			Title:       "Tech Project 1",
			Description: "Test tech project",
			ProjectType: string(models.ProjectTypeTech),
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Title:       "Tech Project 2",
			Description: "Another tech project",
			ProjectType: string(models.ProjectTypeTech),
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	mockDB.On("GetProjects", models.ProjectTypeTech).Return(expectedProjects, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tech/projects", nil)
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
	assert.Equal(t, "Tech Project 1", firstProject["title"])
	assert.Equal(t, "Test tech project", firstProject["description"])
	assert.Equal(t, string(models.ProjectTypeTech), firstProject["project_type"])

	mockDB.AssertExpectations(t)
}

func TestTechController_HandleProjects_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestTechController()

	// Set up mock expectation for empty result
	mockDB.On("GetProjects", models.ProjectTypeTech).Return([]*models.ProjectGroups{}, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tech/projects", nil)
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

func TestTechController_HandleProjects_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestTechController()

	// Set up mock expectation for not found
	mockDB.On("GetProjects", models.ProjectTypeTech).Return(nil, gorm.ErrRecordNotFound)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tech/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockDB.AssertExpectations(t)
}

func TestTechController_HandleProjects_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestTechController()

	// Set up mock expectation for database error
	mockDB.On("GetProjects", models.ProjectTypeTech).Return(nil, errors.New("database connection error"))

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tech/projects", nil)
	router.ServeHTTP(w, req)

	// Should return 500 for database errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockDB.AssertExpectations(t)
}

func TestTechController_HandleProjects_Timeout(t *testing.T) {
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

	controller := NewProjectsController(mockDB, config)

	// Set up mock expectation
	mockDB.On("GetProjects", models.ProjectTypeTech).Return([]*models.ProjectGroups{}, nil).Maybe()

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tech/projects", nil)
	router.ServeHTTP(w, req)

	// Should timeout and return error
	assert.Equal(t, http.StatusRequestTimeout, w.Code)
}
