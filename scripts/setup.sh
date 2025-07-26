#!/bin/bash

set -e

echo "🚀 Setting up Video Agent Go project..."

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p uploads/{images,audio,videos,subtitles}
mkdir -p temp
mkdir -p logs
mkdir -p bin

# Copy environment file
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "Please edit .env file to add your OpenAI API key"
else
    echo "✅ .env file already exists"
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.22 or later"
    exit 1
fi

echo "📦 Downloading Go dependencies..."
go mod download || echo "⚠️  Failed to download some dependencies, continuing..."

# Check if Docker is installed
if command -v docker &> /dev/null; then
    echo "🐳 Docker found"
    
    # Check if docker-compose is available
    if command -v docker-compose &> /dev/null; then
        echo "📋 docker-compose found"
    elif docker compose version &> /dev/null; then
        echo "📋 docker compose (v2) found"
    else
        echo "⚠️  docker-compose not found, please install it for full development experience"
    fi
else
    echo "⚠️  Docker not found, please install it for full development experience"
fi

# Check if FFmpeg is installed
if command -v ffmpeg &> /dev/null; then
    echo "🎬 FFmpeg found"
else
    echo "⚠️  FFmpeg not found. Video processing will not work without it."
    echo "   Install FFmpeg with:"
    echo "   - macOS: brew install ffmpeg"
    echo "   - Ubuntu: sudo apt install ffmpeg"
    echo "   - Or use Docker deployment"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file with your OpenAI API key"
echo "2. Start development: make dev"
echo "3. Or use Docker: make docker-run"
echo ""
echo "📚 See README.md for more details" 