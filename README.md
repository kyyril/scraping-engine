# Scraping Engine: Ultimate CLI & API Scraper

Professional, enterprise-grade web scraping engine built with **Go**. It works both as a dead-simple **CLI tool** for quick scrapes and as a powerful **Distributed API Service** for large-scale automation.

## Dual-Mode Functionality

1.  **CLI Tool**: Perfect for developers. No setup, no database, just run and get data.
2.  **API Service**: Scalable REST API with a job queue, PostgreSQL persistence, and worker pools.

---

## Quick Start (CLI Mode) - No Setup Needed!

Just download/clone and run. No database or Docker required.

### 1. Simple Scrape (Get all text)
```bash
go run cmd/cli/main.go --url "https://example.com"
```

### 2. Extract Specific Data
```bash
go run cmd/cli/main.go --url "https://example.com" --extract "h1, .price, p"
```

### 3. Build into a Binary
```bash
go build -o scrap.exe ./cmd/cli/main.go
./scrap --url "https://google.com"
```

---

## API Service Mode (Docker)

Ideal for production, background jobs, and distributed environments.

### 1. Start with Docker Compose
```bash
docker-compose up -d
```

### 2. Submit a Job (Simplified API)
You can now submit a job with **just a URL**. The system will automatically navigate and extract the page content for you.

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{ "url": "https://example.com" }'
```

### 3. Check Results
```bash
curl http://localhost:8080/api/v1/jobs/{job-id}/result
```

---

## Key Features

- **Headless Chrome**: Uses real browser rendering via `chromedp` (handles SPA/React/Vue).
- **Anti-Detection**: Built-in User-Agent rotation and human-like interaction.
- **Job Queue**: Distributed worker pool handles massive job loads.
- **Persistence**: Auto-saves everything to PostgreSQL (GORM).
- **Swagger UI**: Interactive API testing at `http://localhost:8080/swagger/`.
- **Flexible Actions**: Custom flows (navigate -> wait -> click -> type -> extract).

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI / API     │───▶│   Job Queue     │───▶│ Browser Manager │
│                 │    │                 │    │                 │
│ • Easy Commands │    │ • Worker Pool   │    │ • chromedp      │
│ • JSON Output   │    │ • Retry Logic   │    │ • Session Pool  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    PostgreSQL (Persistence)                     │
└─────────────────────────────────────────────────────────────────┘
```

## Supported Actions (API)

| Action | Description | Parameters |
|--------|-------------|------------|
| `navigate` | Open a URL | `target`: URL |
| `click` | Click an element | `target`: CSS selector |
| `type` | Input text | `target`: CSS selector, `value`: text |
| `wait` | Pause | `value`: seconds |
| `extract` | Get text | `target`: CSS selector, `value`: key name |
| `screenshot`| Take photo | - |
| `scroll` | Scroll down | `options`: `{"x": 0, "y": 500}` |

## Deployment
- **Docker**: `docker build -t scraper .`

## 📄 License
MIT License. Created with ❤️ by Kyyril.
