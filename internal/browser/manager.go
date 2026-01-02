package browser

import (
	"context"
	"distributed-scraper/internal/scraper"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/121.0",
}

type Manager struct {
	maxConcurrent int
	timeout       time.Duration
	activeSessions int
	mutex         sync.Mutex
}

func NewManager(maxConcurrent int, timeout time.Duration) scraper.BrowserManager {
	return &Manager{
		maxConcurrent: maxConcurrent,
		timeout:       timeout,
	}
}

func (m *Manager) CreateSession(ctx context.Context) (scraper.BrowserSession, error) {
	m.mutex.Lock()
	if m.activeSessions >= m.maxConcurrent {
		m.mutex.Unlock()
		return nil, fmt.Errorf("maximum concurrent sessions reached (%d)", m.maxConcurrent)
	}
	m.activeSessions++
	m.mutex.Unlock()

	// Random user agent
	userAgent := userAgents[rand.Intn(len(userAgents))]

	// Create Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.UserAgent(userAgent),
		chromedp.WindowSize(1920, 1080),
	)

	// Create context with timeout
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, m.timeout)

	session := &Session{
		ctx:           timeoutCtx,
		allocCancel:   allocCancel,
		browserCancel: browserCancel,
		timeoutCancel: timeoutCancel,
		manager:       m,
	}

	// Test browser by navigating to a simple page
	if err := chromedp.Run(timeoutCtx, chromedp.Navigate("about:blank")); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to initialize browser session: %w", err)
	}

	return session, nil
}

func (m *Manager) GetAvailableSlots() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.maxConcurrent - m.activeSessions
}

func (m *Manager) Cleanup() {
	log.Println("Cleaning up browser manager...")
	// Additional cleanup logic if needed
}

func (m *Manager) releaseSession() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.activeSessions > 0 {
		m.activeSessions--
	}
}