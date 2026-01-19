# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Loopr is a CLI tool for autonomous development orchestration using Claude Code and the Ralph Loop methodology. It runs Claude in iterative loops to plan and build projects, with each iteration automatically committed to git. The tool is written in Go and uses Docker sandbox for secure Claude execution.

## Tech Stack

- **Language**: Go 1.24+
- **CLI Framework**: urfave/cli v2
- **TUI/UI**: Charmbracelet ecosystem
  - Bubble Tea - Interactive TUI framework
  - Lipgloss - Terminal styling
  - Huh - Form components
- **Sandbox**: Docker (via `docker sandbox run` command)

## Development Commands

### Building
```bash
go build -o loopr main.go
```

### Testing
```bash
go test ./...
```

### Running Locally
```bash
go run main.go <command>
# or after building
./loopr <command>
```

### Cross-Compiling
```bash
# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o loopr-darwin-arm64

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o loopr-darwin-amd64

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o loopr-linux-arm64

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o loopr-linux-amd64
```

### Creating a Release
```bash
./scripts/release.sh
```

The release script will:
1. Check for uncommitted changes
2. Prompt for version (e.g., v0.1.0, v1.0.0-beta.1)
3. Create an annotated git tag
4. Push to GitHub (triggers automated build/release via GitHub Actions)

## Architecture

### Directory Structure

```
loopr/
├── main.go                 # CLI entry point, version handling, command routing
├── cmd/                    # Command implementations
│   ├── init.go            # Initialize .loopr/ directory with templates
│   ├── plan.go            # Run planning loop (also contains Build())
│   └── update.go          # Self-update functionality
├── internal/
│   ├── config/            # Configuration loading from .loopr/config.json
│   ├── sandbox/           # Docker sandbox abstraction and Claude execution
│   ├── git/               # Git operations (push, branch detection)
│   ├── loop/              # Loop controller (iteration management)
│   └── ui/                # Bubble Tea/Huh prompts and styling
├── templates/             # Embedded templates (PROMPT_*.md, AGENTS.md, etc.)
│   └── embed.go          # Go embed directives for templates
└── schema/               # JSON schemas for config validation
```

### Key Architectural Patterns

#### 1. Command Flow
- `main.go` defines CLI commands using urfave/cli
- Commands delegate to functions in `cmd/` package
- `cmd/plan.go` contains both `Plan()` and `Build()` which call shared `runLoop()` function
- `runLoop()` validates environment, loads config, creates sandbox, then delegates to loop controller

#### 2. Loop Controller Pattern
- `internal/loop/controller.go` orchestrates iterations
- Loads prompt file based on mode (plan/build)
- Executes Claude via sandbox for each iteration
- Automatically pushes to git after each iteration
- Handles graceful shutdown on Ctrl+C

#### 3. Sandbox Abstraction
- `internal/sandbox/sandbox.go` defines interface
- `internal/sandbox/docker.go` implements Docker-specific logic
- Uses `docker sandbox run` command with working directory mount
- Authentication check via quick haiku test
- Streams stdout/stderr directly to terminal for live output

#### 4. Template System
- Templates embedded at compile time via `//go:embed`
- `templates/embed.go` provides extraction to `.loopr/` directory
- Templates include: PROMPT_plan.md, PROMPT_build.md, AGENTS.md, config.json, specs/

#### 5. Configuration
- JSON-based configuration in `.loopr/config.json`
- Schema validation (references `schema/config.schema.json`)
- Supports custom models per phase (plan/build)
- Currently only "docker" sandbox is supported

## Important Implementation Details

### Version Handling
- Version is set via ldflags during release builds: `-X main.version=v0.1.0 -X main.commit=abc1234`
- Dev builds show "dev" with git commit from build info
- `getVersion()` in main.go handles both cases

### Docker Sandbox
- Requires Docker Desktop installed and running
- Claude authentication checked via quick haiku test before each command
- Executes with `--dangerously-skip-permissions` flag (user is prompted during init)
- Working directory mounted to allow Claude to modify local files
- Model defaults to "sonnet" unless overridden in config

### Git Operations
- Automatic push after each iteration
- Uses `-u` flag fallback if remote branch doesn't exist
- Requires git repository (checked at start of plan/build)

### Template Extraction
- `loopr init` extracts templates from embedded FS to `.loopr/`
- CLAUDE.md template is extracted to project root (not `.loopr/`)
- Prompts user before overwriting existing `.loopr/` directory

## Code Conventions

- Errors are returned, not logged (except in main.go)
- UI output uses lipgloss styles (defined in `internal/ui/styles.go`)
- All user interaction goes through `internal/ui/prompts.go`
- Git operations isolated in `internal/git/operations.go`
- Sandbox interface allows future non-Docker implementations
