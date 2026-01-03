# Distributed Web Scraping Engine

Enterprise-grade distributed web scraping service with headless browser support, built with Go, Fiber, GORM, and chromedp.

## Features

- **RESTful API** with Fiber for job submission and management
- **Job Queue System** with configurable concurrency limits
- **Headless Browser Integration** using chromedp (Chrome DevTools Protocol)
- **Persistent Storage** with GORM and PostgreSQL
- **Retry Mechanisms** with exponential backoff
- **User-Agent Rotation** for anti-detection
- **Resource Management** with proper cleanup
- **Docker Support** with optimized Alpine + Chromium image
- **Swagger Documentation** for API endpoints
- **Health Checks** and monitoring endpoints

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Fiber API     │───▶│   Job Queue     │───▶│ Browser Manager │
│                 │    │                 │    │                 │
│ • Submit Jobs   │    │ • Worker Pool   │    │ • chromedp      │
│ • Get Results   │    │ • Retry Logic   │    │ • Session Pool  │
│ • Health Check  │    │ • Concurrency   │    │ • Cleanup       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                          │
│                                                                 │
│ • Jobs: Status, Actions, Results                               │
│ • Results: Extracted Data, Screenshots                         │
│ • Metadata: Execution Time, User Agent                        │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd distributed-scraper
```

2. Start the services:
```bash
docker-compose up -d
```

3. The API will be available at `http://localhost:8080`

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Setup PostgreSQL database:
```bash
createdb scraper_db
```

3. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/scraper_db?sslmode=disable"
export PORT=8080
```

4. Run the application:
```bash
go run cmd/main.go
```

## API Usage

### Submit a Scraping Job

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "actions": [
      {
        "type": "navigate",
        "target": "https://example.com"
      },
      {
        "type": "wait",
        "value": "2"
      },
      {
        "type": "extract",
        "target": "h1",
        "value": "page_title"
      },
      {
        "type": "screenshot"
      }
    ],
    "timeout": 30,
    "max_retries": 3
  }'
```

### Get Job Status

```bash
curl http://localhost:8080/api/v1/jobs/{job-id}
```

### Get Job Results

```bash
curl http://localhost:8080/api/v1/jobs/{job-id}/result
```

### Check Queue Status

```bash
curl http://localhost:8080/api/v1/queue/status
```

## Supported Actions

| Action | Description | Parameters |
|--------|-------------|------------|
| `navigate` | Navigate to URL | `target`: URL |
| `click` | Click element | `target`: CSS selector |
| `type` | Type text into element | `target`: CSS selector, `value`: text |
| `wait` | Wait for duration | `value`: seconds |
| `screenshot` | Take screenshot | None |
| `extract` | Extract text from element | `target`: CSS selector, `value`: result key |
| `scroll` | Scroll page | `options`: `{"x": 0, "y": 500}` |

## Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `MAX_CONCURRENT_JOBS`: Maximum concurrent browser sessions (default: 3)
- `BROWSER_TIMEOUT_SECONDS`: Browser session timeout (default: 30)
- `USER_AGENT_ROTATION`: Enable random user agents (default: true)

## Production Deployment

### Railway Deployment

1. Connect your GitHub repository to Railway
2. Add environment variables in Railway dashboard
3. Deploy automatically on git push

### Manual Docker Deployment

```bash
# Build image
docker build -t distributed-scraper .

# Run container
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL="your-postgres-url" \
  -e MAX_CONCURRENT_JOBS=5 \
  distributed-scraper
```

## API Documentation

Visit `http://localhost:8080/swagger/` for interactive Swagger documentation.

## Security Features

- Non-root Docker user for container security
- Input validation and sanitization
- Resource limits to prevent DoS
- Proper error handling without information leakage
- Session isolation and cleanup

## Monitoring

- Health check endpoint: `/health`
- Queue status monitoring: `/api/v1/queue/status`
- Structured logging with request IDs
- Database connection health checks

## Performance Considerations

- Configurable concurrency limits
- Browser session pooling
- Automatic cleanup of stale sessions
- Efficient database queries with proper indexing
- Memory usage monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.