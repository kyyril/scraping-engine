package scraper

import (
	"context"
	"distributed-scraper/internal/models"
	"encoding/base64"
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	browserManager BrowserManager
	db             *gorm.DB
}

func NewService(browserManager BrowserManager, db *gorm.DB) ScraperService {
	return &Service{
		browserManager: browserManager,
		db:             db,
	}
}

func (s *Service) ExecuteJob(ctx context.Context, job *models.ScrapingJob) error {
	// Update job status to processing
	job.Status = models.StatusProcessing
	if err := s.db.Save(job).Error; err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Create browser session
	session, err := s.browserManager.CreateSession(ctx)
	if err != nil {
		return s.markJobFailed(job, fmt.Errorf("failed to create browser session: %w", err))
	}
	defer session.Close()

	// Sort actions by order
	sort.Slice(job.Actions, func(i, j int) bool {
		return job.Actions[i].Order < job.Actions[j].Order
	})

	// Execute actions
	result := &models.ScrapingResult{
		JobID:    job.ID,
		Data:     make(map[string]interface{}),
		Metadata: make(map[string]interface{}),
	}

	for _, action := range job.Actions {
		if err := s.executeAction(session, action, result); err != nil {
			return s.markJobFailed(job, fmt.Errorf("failed to execute action %s: %w", action.Type, err))
		}
	}

	// Save result
	result.Metadata["execution_time"] = time.Since(job.UpdatedAt).Seconds()
	result.Metadata["user_agent"] = job.UserAgent

	if err := s.db.Create(result).Error; err != nil {
		return s.markJobFailed(job, fmt.Errorf("failed to save result: %w", err))
	}

	// Mark job as completed
	now := time.Now()
	job.Status = models.StatusCompleted
	job.CompletedAt = &now
	if err := s.db.Save(job).Error; err != nil {
		log.Printf("Failed to update job completion status: %v", err)
	}

	return nil
}

func (s *Service) executeAction(session BrowserSession, action models.JobAction, result *models.ScrapingResult) error {
	switch action.Type {
	case models.ActionNavigate:
		return session.Navigate(action.Target)
	
	case models.ActionClick:
		return session.Click(action.Target)
	
	case models.ActionType_:
		return session.Type(action.Target, action.Value)
	
	case models.ActionWait:
		duration := 1 // default 1 second
		if action.Value != "" {
			if d, err := time.ParseDuration(action.Value + "s"); err == nil {
				duration = int(d.Seconds())
			}
		}
		return session.Wait(duration)
	
	case models.ActionScreenshot:
		screenshot, err := session.Screenshot()
		if err != nil {
			return err
		}
		screenshotData := base64.StdEncoding.EncodeToString(screenshot)
		result.Screenshots = append(result.Screenshots, screenshotData)
		return nil
	
	case models.ActionExtract:
		text, err := session.ExtractText(action.Target)
		if err != nil {
			return err
		}
		key := action.Value
		if key == "" {
			key = fmt.Sprintf("extracted_%s", action.Target)
		}
		result.Data[key] = text
		return nil
	
	case models.ActionScroll:
		x, y := 0, 500 // default scroll down
		if action.Options != nil {
			if xVal, ok := action.Options["x"].(float64); ok {
				x = int(xVal)
			}
			if yVal, ok := action.Options["y"].(float64); ok {
				y = int(yVal)
			}
		}
		return session.Scroll(x, y)
	
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (s *Service) markJobFailed(job *models.ScrapingJob, err error) error {
	job.Status = models.StatusFailed
	job.Error = err.Error()
	now := time.Now()
	job.CompletedAt = &now
	
	if dbErr := s.db.Save(job).Error; dbErr != nil {
		log.Printf("Failed to update job failure status: %v", dbErr)
	}
	
	return err
}

func (s *Service) GetJob(jobID string) (*models.ScrapingJob, error) {
	var job models.ScrapingJob
	if err := s.db.Preload("Actions").Preload("Result").First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *Service) GetJobResult(jobID string) (*models.ScrapingResult, error) {
	var result models.ScrapingResult
	if err := s.db.First(&result, "job_id = ?", jobID).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) ListJobs(status models.JobStatus, limit, offset int) ([]*models.ScrapingJob, error) {
	var jobs []*models.ScrapingJob
	query := s.db.Preload("Actions").Preload("Result")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, err
	}
	
	return jobs, nil
}