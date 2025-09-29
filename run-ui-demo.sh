#!/bin/bash

echo "🦥 Sloth Runner UI Demo"
echo "======================="

# Build the application
echo "📦 Building Sloth Runner..."
go build -o sloth-runner ./cmd/sloth-runner

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build completed successfully!"

# Start the UI server
echo "🚀 Starting UI server on port 8080..."
echo "📱 Open your browser and go to: http://localhost:8080"
echo ""
echo "🎯 Features available:"
echo "   • 📊 Real-time dashboard"
echo "   • ➕ Create new tasks (Shell, Lua, Pipeline)"
echo "   • 🖥️  Manage agents"
echo "   • 🔄 Live updates via WebSocket"
echo "   • 🌙 Dark/Light theme toggle"
echo "   • 📱 Responsive design"
echo ""
echo "🛑 Press Ctrl+C to stop the server"
echo ""

./sloth-runner ui --port 8080