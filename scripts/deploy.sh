#!/bin/bash

# Deployment script for production

set -e

echo "🚀 Deploying Distributed Scraping Engine..."

# Build Docker image
echo "🐳 Building Docker image..."
docker build -t distributed-scraper:latest .

# Tag for production
docker tag distributed-scraper:latest distributed-scraper:$(date +%Y%m%d-%H%M%S)

# Deploy with Docker Compose
echo "📦 Deploying with Docker Compose..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Health check
echo "🏥 Performing health check..."
sleep 10

if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Deployment successful! Service is healthy."
else
    echo "❌ Deployment failed! Service is not responding."
    exit 1
fi

echo "🎉 Deployment completed successfully!"