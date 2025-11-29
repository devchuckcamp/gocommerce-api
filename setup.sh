#!/bin/bash

# E-Commerce API Setup Script

echo "🚀 Setting up E-Commerce API..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.23 or later."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your configuration before running the app"
else
    echo "✅ .env file already exists"
fi

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
if [ $? -eq 0 ]; then
    echo "✅ Dependencies downloaded successfully"
else
    echo "❌ Failed to download dependencies"
    exit 1
fi

# Tidy up go.mod
echo "🧹 Tidying up dependencies..."
go mod tidy

echo ""
echo "✨ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file with your database credentials"
echo "2. Generate a secure JWT_SECRET (at least 32 characters)"
echo "3. Create your database (e.g., CREATE DATABASE gocommerce;)"
echo "4. Run the application: go run cmd/api/main.go"
echo ""
echo "Optional: Set SEED_DB=true in .env to seed sample data"
echo ""
echo "API will be available at: http://localhost:8080"
echo "Health check: http://localhost:8080/health"
