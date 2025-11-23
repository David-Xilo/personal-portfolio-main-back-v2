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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetContact() (*models.Contacts, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contacts), args.Error(1)
}

func (m *MockDatabase) GetProjects(projectType models.ProjectType) ([]*models.ProjectGroups, error) {
	args := m.Called(projectType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ProjectGroups), args.Error(1)
}

func setupTestAboutController() (*ContactController, *MockDatabase) {
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

	controller := NewContactController(mockDB, config)
	return controller, mockDB
}

func TestNewAboutController(t *testing.T) {
	controller, mockDB := setupTestAboutController()

	assert.NotNil(t, controller)
	assert.Equal(t, mockDB, controller.db)
	assert.NotNil(t, controller.config)
}

func TestAboutController_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, _ := setupTestAboutController()
	router := gin.New()

	controller.RegisterRoutes(router)

	// Test that routes are registered by making requests
	routes := router.Routes()

	// Check that the routes exist
	contactRouteFound := false
	reviewsRouteFound := false

	for _, route := range routes {
		if route.Path == "/about/contact" && route.Method == "GET" {
			contactRouteFound = true
		}
		if route.Path == "/about/reviews/carousel" && route.Method == "GET" {
			reviewsRouteFound = true
		}
	}

	assert.True(t, contactRouteFound, "Contact route should be registered")
	assert.True(t, reviewsRouteFound, "Reviews carousel route should be registered")
}

func TestAboutController_HandleContactRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestAboutController()

	// Set up mock expectation
	expectedContact := &models.Contacts{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Linkedin: "linkedin.com/in/johndoe",
		Github:   "github.com/johndoe",
		Credly:   "credly.com/johndoe",
	}

	mockDB.On("GetContact").Return(expectedContact, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/about/contact", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"].(map[string]interface{})
	assert.Equal(t, "John Doe", message["name"])
	assert.Equal(t, "john@example.com", message["email"])
	assert.Equal(t, "linkedin.com/in/johndoe", message["linkedin"])
	assert.Equal(t, "github.com/johndoe", message["github"])
	assert.Equal(t, "credly.com/johndoe", message["credly"])

	mockDB.AssertExpectations(t)
}

func TestAboutController_HandleContactRequest_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestAboutController()

	// Set up mock expectation for not found
	mockDB.On("GetContact").Return(nil, gorm.ErrRecordNotFound)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/about/contact", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockDB.AssertExpectations(t)
}

func TestAboutController_HandleContactRequest_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestAboutController()

	// Set up mock expectation for database error
	mockDB.On("GetContact").Return(nil, errors.New("database connection error"))

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/about/contact", nil)
	router.ServeHTTP(w, req)

	// Should return 500 for database errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockDB.AssertExpectations(t)
}

func TestAboutController_HandleReviewsCarouselRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, _ := setupTestAboutController()

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/about/reviews/carousel", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	// The message should be an array of reviews
	message := response["message"]
	assert.NotNil(t, message)
}
