package queue

import (
	"context"
	"distributed-scraper/internal/models"
	"distributed-scraper/internal/scraper"
	"fmt"
	"log"
	"sync"
	"time"
)

type JobQueue struct {
	scraperService scraper.ScraperService
	jobs           chan *models.ScrapingJob
	workers        int
	wg             sync.WaitGroup
	stopChan       chan struct{}
	running        bool
	mutex          sync.Mutex
}

func NewJobQueue(scraperService scraper.ScraperService, workers int) *JobQueue {
	return &JobQueue{
		scraperService: scraperService,
		jobs:           make(chan *models.ScrapingJob, workers*2), // Buffer for pending jobs
		workers:        workers,
		stopChan:       make(chan struct{}),
	}
}

func (jq *JobQueue) Start() {
	jq.mutex.Lock()
	defer jq.mutex.Unlock()

	if jq.running {
		return
	}

	jq.running = true
	log.Printf("Starting job queue with %d workers", jq.workers)

	// Start worker goroutines
	for i := 0; i < jq.workers; i++ {
		jq.wg.Add(1)
		go jq.worker(i)
	}
}

func (jq *JobQueue) Stop() {
	jq.mutex.Lock()
	defer jq.mutex.Unlock()

	if !jq.running {
		return
	}

	log.Println("Stopping job queue...")
	jq.running = false
	close(jq.stopChan)
	jq.wg.Wait()
	log.Println("Job queue stopped")
}

func (jq *JobQueue) SubmitJob(job *models.ScrapingJob) error {
	jq.mutex.Lock()
	defer jq.mutex.Unlock()

	if !jq.running {
		return fmt.Errorf("job queue is not running")
	}

	select {
	case jq.jobs <- job:
		log.Printf("Job %s submitted to queue", job.ID)
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

func (jq *JobQueue) worker(id int) {
	defer jq.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-jq.stopChan:
			log.Printf("Worker %d stopping", id)
			return
		
		case job := <-jq.jobs:
			if job == nil {
				continue
			}

			log.Printf("Worker %d processing job %s", id, job.ID)
			
			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
			
			// Execute job with retry logic
			err := jq.executeJobWithRetry(ctx, job)
			if err != nil {
				log.Printf("Worker %d failed to process job %s: %v", id, job.ID, err)
			} else {
				log.Printf("Worker %d completed job %s", id, job.ID)
			}
			
			cancel()
		}
	}
}

func (jq *JobQueue) executeJobWithRetry(ctx context.Context, job *models.ScrapingJob) error {
	var lastError error

	for attempt := 0; attempt <= job.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := time.Duration(attempt*attempt) * time.Second
			log.Printf("Retrying job %s in %v (attempt %d/%d)", job.ID, delay, attempt, job.MaxRetries)
			time.Sleep(delay)
		}

		err := jq.scraperService.ExecuteJob(ctx, job)
		if err == nil {
			return nil
		}

		lastError = err
		job.Retries = attempt + 1

		// Check if we should stop retrying
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return lastError
}

func (jq *JobQueue) GetQueueSize() int {
	return len(jq.jobs)
}

func (jq *JobQueue) IsRunning() bool {
	jq.mutex.Lock()
	defer jq.mutex.Unlock()
	return jq.running
}