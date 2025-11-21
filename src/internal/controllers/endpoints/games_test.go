package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	configuration "personal-portfolio-main-back/src/internal/config"
	"personal-portfolio-main-back/src/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestGamesController() (*GamesController, *MockDatabase) {
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

	controller := NewGamesController(mockDB, config)
	return controller, mockDB
}

func TestNewGamesController(t *testing.T) {
	controller, mockDB := setupTestGamesController()

	assert.NotNil(t, controller)
	assert.Equal(t, mockDB, controller.db)
	assert.NotNil(t, controller.config)
}

func TestGamesController_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, _ := setupTestGamesController()
	router := gin.New()

	controller.RegisterRoutes(router)

	// Test that routes are registered by checking the routes
	routes := router.Routes()

	// Check that both routes exist
	projectsRouteFound := false
	carouselRouteFound := false

	for _, route := range routes {
		if route.Path == "/games/projects" && route.Method == "GET" {
			projectsRouteFound = true
		}
		if route.Path == "/games/played/carousel" && route.Method == "GET" {
			carouselRouteFound = true
		}
	}

	assert.True(t, projectsRouteFound, "Games projects route should be registered")
	assert.True(t, carouselRouteFound, "Games played carousel route should be registered")
}

func TestGamesController_HandleProjects_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestGamesController()

	// Set up mock expectation
	expectedProjects := []*models.ProjectGroups{
		{
			ID:          1,
			Title:       "Game Project 1",
			Description: "Test game project",
			ProjectType: string(models.ProjectTypeGame),
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	mockDB.On("GetProjects", models.ProjectTypeGame).Return(expectedProjects, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"].([]interface{})
	assert.Len(t, message, 1)

	// Check project
	project := message[0].(map[string]interface{})
	assert.Equal(t, "Game Project 1", project["title"])
	assert.Equal(t, "Test game project", project["description"])
	assert.Equal(t, string(models.ProjectTypeGame), project["project_type"])

	mockDB.AssertExpectations(t)
}

func TestGamesController_HandleProjects_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestGamesController()

	// Set up mock expectation for database error
	mockDB.On("GetProjects", models.ProjectTypeGame).Return(nil, errors.New("database connection error"))

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/projects", nil)
	router.ServeHTTP(w, req)

	// Should return 500 for database errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockDB.AssertExpectations(t)
}

func TestGamesController_HandleGamesPlayedCarousel_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestGamesController()

	// Set up mock expectation
	expectedGames := []*models.GamesPlayed{
		{
			ID:          1,
			Title:       "Game 1",
			Description: "First game",
			Rating:      5,
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Title:       "Game 2",
			Description: "Second game",
			Rating:      4,
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	mockDB.On("GetGamesPlayed").Return(expectedGames, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/played/carousel", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"].([]interface{})
	assert.Len(t, message, 2)

	// Check first game
	firstGame := message[0].(map[string]interface{})
	assert.Equal(t, "Game 1", firstGame["title"])
	assert.Equal(t, "First game", firstGame["description"])
	assert.Equal(t, float64(5), firstGame["rating"])

	mockDB.AssertExpectations(t)
}

func TestGamesController_HandleGamesPlayedCarousel_LimitToFive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestGamesController()

	// Set up mock expectation with more than 5 games
	expectedGames := make([]*models.GamesPlayed, 7)
	for i := 0; i < 7; i++ {
		expectedGames[i] = &models.GamesPlayed{
			ID:          uint(i + 1),
			Title:       fmt.Sprintf("Game %d", i+1),
			Description: fmt.Sprintf("Game description %d", i+1),
			Rating:      5,
			CreatedAt:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		}
	}

	mockDB.On("GetGamesPlayed").Return(expectedGames, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/played/carousel", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"].([]interface{})
	// Should be limited to 5 even though we returned 7
	assert.Len(t, message, 5)

	mockDB.AssertExpectations(t)
}

func TestGamesController_HandleGamesPlayedCarousel_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestGamesController()

	// Set up mock expectation for empty result
	mockDB.On("GetGamesPlayed").Return([]*models.GamesPlayed{}, nil)

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/played/carousel", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")

	message := response["message"]
	if message != nil {
		messageSlice := message.([]interface{})
		assert.Len(t, messageSlice, 0)
	} else {
		// If message is nil, that's also acceptable for empty result
		assert.Nil(t, message)
	}

	mockDB.AssertExpectations(t)
}

func TestGamesController_HandleGamesPlayedCarousel_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller, mockDB := setupTestGamesController()

	// Set up mock expectation for database error
	mockDB.On("GetGamesPlayed").Return(nil, errors.New("database connection error"))

	// Create test request
	router := gin.New()
	controller.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/played/carousel", nil)
	router.ServeHTTP(w, req)

	// Should return 500 for database errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockDB.AssertExpectations(t)
}

func TestGamesFilter_Struct(t *testing.T) {
	filter := GamesFilter{
		Genre: "Action",
	}

	assert.Equal(t, "Action", filter.Genre)
}
