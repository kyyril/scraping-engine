package main

import (
	"distributed-scraper/internal/api"
	_ "distributed-scraper/docs"
	"distributed-scraper/internal/browser"
	"distributed-scraper/internal/middleware"
	"distributed-scraper/internal/models"
	"distributed-scraper/internal/queue"
	"distributed-scraper/internal/scraper"
	"distributed-scraper/pkg/config"
	"distributed-scraper/pkg/database"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title Distributed Web Scraper API
// @version 1.0
// @description Enterprise-grade distributed web scraping service with headless browser support
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&models.ScrapingJob{}, &models.ScrapingResult{}, &models.JobAction{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize browser manager
	browserManager := browser.NewManager(cfg.MaxConcurrentJobs, cfg.BrowserTimeout)

	// Initialize scraper service
	scraperService := scraper.NewService(browserManager, db)

	// Initialize job queue
	jobQueue := queue.NewJobQueue(scraperService, cfg.MaxConcurrentJobs)
	jobQueue.Start()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.RequestLogger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	
	// API Key authentication (if configured)
	if cfg.APIKey != "" {
		app.Use("/api", middleware.APIKeyAuth(cfg.APIKey))
	}

	// Setup API routes
	api.SetupRoutes(app, scraperService, jobQueue, db)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		jobQueue.Stop()
		browserManager.Cleanup()
		app.Shutdown()
	}()

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}