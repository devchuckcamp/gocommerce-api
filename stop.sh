#!/bin/bash

# GoShop E-Commerce API - Stop Local Development

echo "🛑 Stopping GoShop E-Commerce API (local)..."
echo "============================================"
echo ""

# Find and kill Go API process
API_PID=$(ps aux | grep "cmd/api/main.go" | grep -v grep | awk '{print $2}')

if [ -z "$API_PID" ]; then
    echo "ℹ️  No running API process found."
    echo ""
    echo "💡 Tip: If you're running the API in a terminal, use Ctrl+C to stop it."
else
    echo "📍 Found API process: PID $API_PID"
    kill $API_PID 2>/dev/null
    
    # Wait a moment and check if it's stopped
    sleep 1
    if ps -p $API_PID > /dev/null 2>&1; then
        echo "⚠️  Process still running, forcing stop..."
        kill -9 $API_PID 2>/dev/null
    fi
    
    echo "✅ API stopped successfully!"
fi

echo ""
echo "📋 Useful Commands:"
echo "   Start again:      go run cmd/api/main.go"
echo "   Or use setup:     ./setup.sh"
echo ""
