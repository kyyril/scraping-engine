package main

import (
	"context"
	"distributed-scraper/internal/browser"
	"distributed-scraper/internal/models"
	"distributed-scraper/internal/scraper"
	"distributed-scraper/pkg/config"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Parse CLI flags
	url := flag.String("url", "", "URL to scrape")
	selector := flag.String("extract", "", "CSS selectors to extract text from (comma-separated)")
	headless := flag.Bool("headless", true, "Run in headless mode (default true)")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: --url is required")
		flag.Usage()
		os.Exit(1)
	}

	// Setup minimal minimal service dependencies
	cfg := config.Load()
	browserManager := browser.NewManager(1, cfg.BrowserTimeout)
	scraperService := scraper.NewService(browserManager, nil)

	// Determine extractors
	selectors := []string{}
	if *selector != "" {
		parts := strings.Split(*selector, ",")
		for _, p := range parts {
			selectors = append(selectors, strings.TrimSpace(p))
		}
	} else {
		// Use body as default to get everything
		selectors = []string{"body"}
	}

	// Construct actions
	actions := []models.JobAction{
		{
			Type:   models.ActionNavigate,
			Target: *url,
			Order:  0,
		},
	}

	for i, s := range selectors {
		actions = append(actions, models.JobAction{
			Type:   models.ActionExtract,
			Target: s,
			Value:  s, // The key in the result map
			Order:  i + 1,
		})
	}

	// Construct a job object
	job := &models.ScrapingJob{
		URL:       *url,
		Actions:   actions,
		UserAgent: "CLI-Scraper/1.0",
		Timeout:   30,
	}

	fmt.Printf("🕷️ Scraping %s...\n", *url)

	// Create a browser session manually
	ctx := context.Background()
	if !*headless {
		// If you wanted to support non-headless, you'd need to modify BrowserManager
		// For now, it respects the manager's defaults which are likely headless
	}

	session, err := browserManager.CreateSession(ctx)
	if err != nil {
		log.Fatalf("Failed to create browser session: %v", err)
	}
	defer session.Close()

	// Execute pure logic
	result, err := scraperService.ProcessJob(ctx, session, job)
	if err != nil {
		log.Fatalf("❌ Scraping failed: %v", err)
	}

	// Output results as JSON
	fmt.Println("✅ Scraping successful!")
	output, _ := json.MarshalIndent(result.Data, "", "  ")
	fmt.Println(string(output))
}
