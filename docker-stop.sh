#!/bin/bash

# GoShop E-Commerce API - Docker Stop Script

set -e

echo "🛑 Stopping GoShop E-Commerce API..."
echo "===================================="
echo ""

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed."
    exit 1
fi

# Check if containers are running
if ! docker-compose ps | grep -q "Up"; then
    echo "ℹ️  No running containers found."
    exit 0
fi

# Show what will be stopped
echo "📦 Current services:"
docker-compose ps
echo ""

# Stop containers
echo "🛑 Stopping containers..."
docker-compose stop

echo ""
echo "✅ All containers stopped successfully!"
echo ""
echo "📋 Useful Commands:"
echo "   Start again:      ./docker-start.sh"
echo "   Remove volumes:   docker-compose down -v"
echo "   Remove all:       docker-compose down -v --rmi all"
echo ""
