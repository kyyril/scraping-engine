FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o docs
RUN CGO_ENABLED=0 GOOS=linux go build -o scraper-api ./cmd/main.go

FROM alpine:latest

# Install Chromium and dependencies for headless browsing
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    && rm -rf /var/cache/apk/*

# Environment variables for chromedp
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

# Create non-root user for security
RUN addgroup -g 1001 -S scraper && \
    adduser -S scraper -u 1001 -G scraper

WORKDIR /app

COPY --from=builder /app/scraper-api .
RUN chmod +x scraper-api

# Change ownership to non-root user
RUN chown -R scraper:scraper /app

USER scraper

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./scraper-api"]