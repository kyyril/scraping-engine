package api

import (
	"distributed-scraper/internal/models"
	"distributed-scraper/internal/queue"
	"distributed-scraper/internal/scraper"
	"distributed-scraper/internal/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, scraperService scraper.ScraperService, jobQueue *queue.JobQueue, db *gorm.DB) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":     "healthy",
			"timestamp":  time.Now(),
			"queue_size": jobQueue.GetQueueSize(),
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := app.Group("/api/v1")

	// Jobs endpoints
	jobs := api.Group("/jobs")
	jobs.Post("/", createJob(scraperService, jobQueue, db))
	jobs.Get("/", listJobs(scraperService))
	jobs.Get("/:id", getJob(scraperService))
	jobs.Get("/:id/result", getJobResult(scraperService))

	// Queue management
	queue := api.Group("/queue")
	queue.Get("/status", getQueueStatus(jobQueue))
}

// @Summary Create a new scraping job
// @Description Submit a new web scraping job to the queue
// @Tags jobs
// @Accept json
// @Produce json
// @Param job body CreateJobRequest true "Job configuration"
// @Success 201 {object} models.ScrapingJob
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/jobs [post]
func createJob(scraperService scraper.ScraperService, jobQueue *queue.JobQueue, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateJobRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Invalid request body: " + err.Error(),
			})
		}

		// Validate request
		if req.URL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "URL is required",
			})
		}

		// Validate URL format and security
		if err := utils.ValidateURL(req.URL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Invalid URL: " + err.Error(),
			})
		}
		if len(req.Actions) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "At least one action is required",
			})
		}

		// Validate actions
		for i, actionReq := range req.Actions {
			if err := utils.ValidateActionType(actionReq.Type); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
					Error: fmt.Sprintf("Action %d: %s", i, err.Error()),
				})
			}
			
			if err := utils.ValidateSelector(actionReq.Type, actionReq.Target); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
					Error: fmt.Sprintf("Action %d: %s", i, err.Error()),
				})
			}
		}
		// Create job
		job := &models.ScrapingJob{
			URL:        req.URL,
			UserAgent:  req.UserAgent,
			Timeout:    req.Timeout,
			MaxRetries: req.MaxRetries,
			Status:     models.StatusPending,
		}

		if job.Timeout == 0 {
			job.Timeout = 30 // default 30 seconds
		}
		if job.MaxRetries == 0 {
			job.MaxRetries = 3 // default 3 retries
		}

		// Create actions
		for i, actionReq := range req.Actions {
			action := models.JobAction{
				Type:    models.ActionType(actionReq.Type),
				Target:  actionReq.Target,
				Value:   actionReq.Value,
				Options: actionReq.Options,
				Order:   i,
			}
			job.Actions = append(job.Actions, action)
		}

		// Save to database
		if err := db.Create(job).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to create job: " + err.Error(),
			})
		}

		// Submit to queue
		if err := jobQueue.SubmitJob(job); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to submit job to queue: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(job)
	}
}

// @Summary List scraping jobs
// @Description Get a list of scraping jobs with optional filtering
// @Tags jobs
// @Produce json
// @Param status query string false "Filter by status" Enums(pending, processing, completed, failed)
// @Param limit query int false "Number of jobs to return" default(10)
// @Param offset query int false "Number of jobs to skip" default(0)
// @Success 200 {array} models.ScrapingJob
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/jobs [get]
func listJobs(scraperService scraper.ScraperService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status := models.JobStatus(c.Query("status", ""))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))

		// Validate limit
		if limit > 100 {
			limit = 100 // Maximum 100 items per page
		}
		if limit < 1 {
			limit = 10
		}
		jobs, err := scraperService.ListJobs(status, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to fetch jobs: " + err.Error(),
			})
		}

		return c.JSON(jobs)
	}
}

// @Summary Get a specific job
// @Description Get details of a specific scraping job by ID
// @Tags jobs
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} models.ScrapingJob
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/jobs/{id} [get]
func getJob(scraperService scraper.ScraperService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("id")
		
		job, err := scraperService.GetJob(jobID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
					Error: "Job not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to fetch job: " + err.Error(),
			})
		}

		return c.JSON(job)
	}
}

// @Summary Get job result
// @Description Get the result of a completed scraping job
// @Tags jobs
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} models.ScrapingResult
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/jobs/{id}/result [get]
func getJobResult(scraperService scraper.ScraperService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("id")
		
		result, err := scraperService.GetJobResult(jobID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
					Error: "Job result not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to fetch job result: " + err.Error(),
			})
		}

		return c.JSON(result)
	}
}

// @Summary Get queue status
// @Description Get current status of the job queue
// @Tags queue
// @Produce json
// @Success 200 {object} QueueStatusResponse
// @Router /api/v1/queue/status [get]
func getQueueStatus(jobQueue *queue.JobQueue) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(QueueStatusResponse{
			QueueSize: jobQueue.GetQueueSize(),
			IsRunning: jobQueue.IsRunning(),
		})
	}
}