#!/bin/bash

echo "🚀 Starting VamsaSetu - Family Tree & Event Management System"
echo "=============================================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

echo "📦 Starting all services with Docker Compose..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

echo "🔍 Checking service status..."
docker-compose ps

echo ""
echo "✅ VamsaSetu is now running!"
echo ""
echo "🌐 Access the application:"
echo "   Frontend:     http://localhost:3000"
echo "   Backend API:  http://localhost:8080/api"
echo "   Neo4j:        http://localhost:7474 (neo4j/vamsasetu123)"
echo "   PostgreSQL:   localhost:5432 (vamsasetu/vamsasetu123)"
echo ""
echo "📚 For more information, see README.md"
echo ""
echo "🛑 To stop the application, run: docker-compose down"
