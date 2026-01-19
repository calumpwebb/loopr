# Claude Instructions

## Project Overview

Loopr is a CLI tool for autonomous development orchestration using Claude Code and the Ralph Loop methodology. It runs Claude in iterative loops to plan and build projects, with each iteration automatically committed to git.

## Tech Stack

- Go 1.22+
- urfave/cli - CLI framework
- Bubble Tea (charmbracelet) - Interactive TUI
- Lipgloss (charmbracelet) - Terminal styling
- Huh (charmbracelet) - Form components
- Docker Sandbox - Secure Claude execution

## Architecture

```
loopr/
├── main.go              # CLI entry point
├── cmd/                 # Command implementations (init, plan, build)
├── internal/            # Internal packages
│   ├── config/         # Configuration loading
│   ├── sandbox/        # Docker sandbox abstraction
│   ├── git/            # Git operations
│   ├── loop/           # Loop controller
│   └── ui/             # Bubble Tea UI components
├── templates/          # Embedded templates
│   ├── PROMPT_plan.md
│   ├── PROMPT_build.md
│   ├── AGENTS.md
│   ├── config.json
│   └── specs/
└── schema/            # JSON schemas
```

## Loopr Integration

This project uses Loopr itself for development (dogfooding).

### Running Loops

```bash
# Initialize loopr
loopr init

# Generate implementation plan
loopr plan

# Build from plan
loopr build
```

## Development Guidelines

- **Go formatting**: Pre-commit hook runs `go fmt ./...` automatically
- **Testing**: Run `go test ./...` before committing
- **Building**: `go build -o loopr` to build locally
- **Cross-compile**: Use GOOS/GOARCH for other platforms

## Code Conventions

- Use Go standard library where possible
- Errors are returned, not logged (except in main)
- UI components use charmbracelet libraries consistently
- All templates are embedded at compile time
- Config is JSON with JSON Schema validation
