# VamsaSetu Quick Start Guide

## Prerequisites

You need to install:
1. **Docker Desktop** - https://www.docker.com/products/docker-desktop/
2. **Go** (for running seed script) - https://go.dev/doc/install

### Installing Go on Mac

```bash
# Using Homebrew (recommended)
brew install go

# Or download from https://go.dev/dl/
```

Verify installation:
```bash
go version
```

## Step 1: Setup Environment Variables

Create a `.env` file in the project root:

```bash
# Copy the example file
cp .env.example .env
```

Edit `.env` and set at minimum:
```
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

The other variables (Twilio, SendGrid) are optional for basic testing.

## Step 2: Start All Services

```bash
# From the project root directory
docker-compose up -d

# Wait 30-60 seconds for all services to start
# Check status
docker-compose ps
```

You should see all services as "healthy" or "running".

## Step 3: Initialize Go Modules (First Time Only)

```bash
cd backend
go mod download
cd ..
```

## Step 4: Load Demo Data

```bash
cd backend
go run cmd/seed/main.go
cd ..
```

This creates:
- Demo user: `demo@vamsasetu.com` / `Demo@1234`
- 12 family members
- Relationships between them
- 5 sample events

## Step 5: Access the Application

Open your browser:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api
- **Health Check**: http://localhost:8080/health

Login with: `demo@vamsasetu.com` / `Demo@1234`

## Troubleshooting

### Services won't start
```bash
# Stop everything
docker-compose down

# Remove old containers and volumes
docker-compose down -v

# Rebuild and start
docker-compose up -d --build
```

### Check logs
```bash
# All services
docker-compose logs

# Specific service
docker-compose logs backend
docker-compose logs frontend
```

### Go command not found
Install Go first:
```bash
brew install go
```

### Port already in use
Stop the conflicting service or change ports in docker-compose.yml

## Development Mode (Without Docker)

If you want to run services locally for development:

### Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

Make sure databases (Neo4j, PostgreSQL, Redis) are still running via Docker:
```bash
docker-compose up -d neo4j postgres redis
```

## Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```
