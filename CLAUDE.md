# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a multi-modal video generation service built with Go and the Hertz framework. The system uses a Tool-based AI Agent architecture where an LLM (GPT-4) dynamically selects and orchestrates tools to generate videos based on user requests.

## Key Commands

### Development
- `make dev` - Start development environment with MySQL/Redis dependencies
- `make dev-watch` - Run with hot reload using air (requires `make dev-tools` first)
- `make test` - Run all tests
- `go test -v ./agent/` - Run tests for specific package

### Build & Deploy
- `make build` - Build the application binary
- `make docker-build` - Build Docker image
- `make docker-run` - Run with Docker Compose
- `make prod-build` - Create production build (CGO disabled, Linux target)

### Code Quality
- `make fmt` - Format all Go code
- `make lint` - Run golangci-lint (requires `make dev-tools` first)

### Setup
- `make setup` - Initial project setup (creates directories, copies .env)
- `make deps` - Download and tidy Go dependencies

## Architecture

### Core Design: Tool-based AI Agent System

The system implements three execution modes:

1. **Fixed Workflow** (`/api/v1/video/generate`) - Hardcoded 4-step process
2. **LLM Agent Orchestration** (`/api/v1/video/generate-smart`) - LLM selects predefined agents
3. **Tool-based Orchestration** (`/api/v1/video/generate-tools`) - LLM dynamically selects and combines tools

### Key Components

- **Tool Registry** (`agent/tools.go`) - Manages available tools with standardized interface
- **Tool Orchestrator** (`agent/tool_orchestrator.go`) - LLM-driven tool selection and execution
- **Sub-Agents** (`agent/sub_agents.go`) - Specialized agents for different content types
- **Tools** - Implement the `Tool` interface:
  - ContentAnalysisTool - Analyzes user input
  - ScriptGenerationTool - Generates video scripts
  - ImageGenerationTool - Creates images via DALL-E 3
  - VoiceGenerationTool - Text-to-speech synthesis
  - QualityCheckTool - Content quality validation
  - VideoRenderTool - FFmpeg-based video rendering

### LLM Decision Flow

1. User request → LLM analyzes with system prompt
2. LLM selects appropriate tools based on context
3. Tools execute and return results
4. LLM evaluates results and decides next steps
5. Process continues until video is generated or max iterations reached

## Configuration

### Required Environment Variables
- `OPENAI_API_KEY` - OpenAI API key (required)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - MySQL configuration
- `SERVER_PORT` - HTTP server port (default: 8080)
- `STORAGE_TYPE` - Storage backend: "local" or "cloud"

### Directory Structure
- `uploads/` - File storage (images, audio, videos, subtitles)
- `temp/` - Temporary processing files
- `logs/` - Application logs

## Dependencies

- **Web Framework**: Hertz (ByteDance's high-performance HTTP framework)
- **Database**: MySQL 8.0+ with go-sql-driver
- **AI Services**: OpenAI GPT-4 and DALL-E 3
- **Media Processing**: FFmpeg (must be installed on system or use Docker)
- **Configuration**: godotenv for environment management

## Development Tips

1. The system uses OpenAI's function calling API for tool orchestration
2. Each tool must implement the `Tool` interface with proper parameter schemas
3. Tool results are passed back to the LLM for decision making
4. The orchestrator maintains context across multiple tool calls
5. Maximum 10 iterations per request to prevent infinite loops

## Database

MySQL schema is initialized via `init.sql`. The main table tracks video generation tasks with status, metadata, and file paths.

## Testing

Currently no unit tests are implemented. When adding tests:
- Place them in the same package with `_test.go` suffix
- Use the standard Go testing framework
- Mock external services (OpenAI, storage) for reliable tests