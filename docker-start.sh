#!/bin/bash

# GoShop E-Commerce API - Docker Quick Start Script

set -e

echo "🚀 GoShop E-Commerce API - Docker Deployment"
echo "=============================================="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
    echo ""
    echo "⚠️  IMPORTANT: Please update the following in .env:"
    echo "   - JWT_SECRET (must be at least 32 characters)"
    echo "   - DB_PASSWORD (recommended for production)"
    echo "   - GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET (if using OAuth)"
    echo ""
    read -p "Press Enter to continue after updating .env, or Ctrl+C to exit..."
fi

# Check if JWT_SECRET is set
JWT_SECRET=$(grep JWT_SECRET .env | cut -d '=' -f2)
if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" = "your-super-secret-jwt-key-minimum-32-characters-long-change-this" ]; then
    echo "⚠️  Warning: JWT_SECRET is not set or using default value"
    echo "🔐 Generating secure JWT secret..."
    NEW_SECRET=$(openssl rand -base64 48 | tr -d '\n')
    sed -i.bak "s|JWT_SECRET=.*|JWT_SECRET=$NEW_SECRET|" .env
    echo "✅ Generated new JWT_SECRET"
    echo ""
fi

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Stop existing containers
echo "🛑 Stopping existing containers..."
docker-compose down 2>/dev/null || true
echo ""

# Build and start services
echo "🏗️  Building and starting services..."
echo "   This may take a few minutes on first run..."
echo ""
docker-compose up --build -d

# Wait for services to be healthy
echo ""
echo "⏳ Waiting for services to be ready..."
echo ""

# Wait for postgres
echo -n "   Waiting for PostgreSQL"
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U commerce -d commerce &>/dev/null; then
        echo " ✅"
        break
    fi
    echo -n "."
    sleep 1
done

# Wait for API
echo -n "   Waiting for API"
for i in {1..30}; do
    if curl -s http://localhost:8080/health &>/dev/null; then
        echo " ✅"
        break
    fi
    echo -n "."
    sleep 1
done

echo ""
echo "=============================================="
echo "✅ GoShop E-Commerce API is ready!"
echo "=============================================="
echo ""
echo "📍 API Endpoints:"
echo "   - Health Check: http://localhost:8080/health"
echo "   - API Base URL: http://localhost:8080/api/v1"
echo "   - Swagger Docs: http://localhost:8080/api/v1/docs (if enabled)"
echo ""
echo "🔐 Authentication:"
echo "   - Register: POST /api/v1/auth/register"
echo "   - Login:    POST /api/v1/auth/login"
echo "   - OAuth:    GET  /api/v1/auth/google"
echo ""
echo "📦 Services:"
docker-compose ps
echo ""
echo "📋 Useful Commands:"
echo "   View logs:        docker-compose logs -f"
echo "   Stop services:    docker-compose down"
echo "   Restart:          docker-compose restart"
echo "   Database shell:   docker-compose exec postgres psql -U commerce -d commerce"
echo ""
echo "📚 Documentation:"
echo "   - Full guide: ./DOCKER.md"
echo "   - API docs:   ./API.md"
echo "   - OAuth:      ./GOOGLE_OAUTH.md"
echo ""
echo "Happy coding! 🎉"
