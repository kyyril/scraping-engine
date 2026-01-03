# Deployment Guide

## 🚀 Deploying to Railway

This project is configured to be easily deployed to [Railway](https://railway.app/).

### Prerequisites

- A GitHub account with this repository pushed.
- A Railway account.

### Steps

1.  **New Project**: In Railway, click "New Project" > "Deploy from GitHub repo".
2.  **Select Repository**: Choose this repository.
3.  **Add Database**:
    - Right-click on the canvas or click "New" > "Database" > "PostgreSQL".
    - This will automatically create a Postgres instance.
4.  **Connect Service to Database**:
    - Railway automatically provides `DATABASE_URL` environment variable if you link them, wait for the variable to populate or add it manually in the "Variables" tab of your service.
    - Go to your Service > Variables.
    - Ensure `DATABASE_URL` is set (Railway often does this automatically if you add the DB plugin, otherwise copy the connection string from the Postgres service).
5.  **Environment Variables**:
    Add the following variables in the "Variables" tab:
    - `PORT`: `8080` (Railway usually sets this, but it's good to be explicit or let the code default).
    - `MAX_CONCURRENT_JOBS`: `3` (Adjust based on your plan's memory limits).
    - `BROWSER_TIMEOUT_SECONDS`: `60`
    - `user_agent_rotation`: `true`
6.  **Build & Deploy**:
    - Railway will detect the `Dockerfile` in the root directory.
    - It will build the image (installing Chrome and Go dependencies) and deploy it.
7.  **Verify**:
    - Check the "Deployments" log to see the build progress.
    - Once active, use the provided public URL (e.g., `https://web-production-xxxx.up.railway.app`) to access the API.
    - Endpoint: `GET /health` to verify status.

### Notes on Memory
Headless Chrome can be memory intensive. On the Railway starter plan, keep `MAX_CONCURRENT_JOBS` low (e.g., 1 or 2) to avoid OOM (Out of Memory) kills.

### Dockerfile Details
The included `Dockerfile` is optimized for Railway:
- Uses `alpine` for a small footprint.
- Installs `chromium` and strict dependencies.
- Runs as a non-root user (`scraper`) for security, which is best practice on Railway.
