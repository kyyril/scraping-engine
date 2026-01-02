package scraper

import (
	"context"
	"distributed-scraper/internal/models"
)

// BrowserManager defines the interface for browser management
type BrowserManager interface {
	CreateSession(ctx context.Context) (BrowserSession, error)
	GetAvailableSlots() int
	Cleanup()
}

// BrowserSession defines the interface for browser sessions
type BrowserSession interface {
	Navigate(url string) error
	Click(selector string) error
	Type(selector, text string) error
	Wait(duration int) error
	Screenshot() ([]byte, error)
	ExtractText(selector string) (string, error)
	ExtractAttribute(selector, attribute string) (string, error)
	Scroll(x, y int) error
	Close() error
}

// ScraperService defines the interface for the main scraping service
type ScraperService interface {
	ExecuteJob(ctx context.Context, job *models.ScrapingJob) error
	GetJob(jobID string) (*models.ScrapingJob, error)
	GetJobResult(jobID string) (*models.ScrapingResult, error)
	ListJobs(status models.JobStatus, limit, offset int) ([]*models.ScrapingJob, error)
}