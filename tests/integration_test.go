package tests

import (
	"bytes"
	"distributed-scraper/internal/models"
	"distributed-scraper/pkg/config"
	"distributed-scraper/pkg/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	app := fiber.New()
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateJobValidation(t *testing.T) {
	app := fiber.New()
	
	app.Post("/api/v1/jobs", func(c *fiber.Ctx) error {
		var req struct {
			URL     string `json:"url"`
			Actions []struct {
				Type   string `json:"type"`
				Target string `json:"target"`
			} `json:"actions"`
		}
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		
		if req.URL == "" {
			return c.Status(400).JSON(fiber.Map{"error": "URL is required"})
		}
		
		if len(req.Actions) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "At least one action is required"})
		}
		
		return c.Status(201).JSON(fiber.Map{"id": "test-job-id"})
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Missing URL",
			payload:        map[string]interface{}{"actions": []interface{}{}},
			expectedStatus: 400,
			expectedError:  "URL is required",
		},
		{
			name:           "Missing actions",
			payload:        map[string]interface{}{"url": "https://example.com"},
			expectedStatus: 400,
			expectedError:  "At least one action is required",
		},
		{
			name: "Valid request",
			payload: map[string]interface{}{
				"url": "https://example.com",
				"actions": []interface{}{
					map[string]interface{}{"type": "navigate", "target": "https://example.com"},
				},
			},
			expectedStatus: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/v1/jobs", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestDatabaseConnection(t *testing.T) {
	cfg := &config.Config{
		DatabaseURL: "postgres://scraper:scraper_password@localhost:5432/scraper_test?sslmode=disable",
	}
	
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		t.Skip("Database not available for testing")
	}
	
	// Test migration
	err = db.AutoMigrate(&models.ScrapingJob{}, &models.ScrapingResult{}, &models.JobAction{})
	assert.NoError(t, err)
	
	// Test creating a job
	job := &models.ScrapingJob{
		URL:     "https://example.com",
		Status:  models.StatusPending,
		Timeout: 30,
	}
	
	err = db.Create(job).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, job.ID)
	
	// Cleanup
	db.Delete(job)
}