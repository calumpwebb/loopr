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

## Testing Loopr with tmux

**IMPORTANT**: Because loopr is highly interactive (prompts, theme selection, authentication), you MUST use tmux when testing programmatically. Direct bash execution will fail on interactive prompts.

### Why tmux?
- Loopr has interactive prompts (iterations, authentication, Claude theme setup)
- tmux allows you to send keystrokes and capture output from interactive sessions
- You can monitor the session in real-time while Claude runs

### Preparing for Tests

**CRITICAL**: Always clear Docker sandboxes before testing authentication flows or new states:

```bash
# List all sandboxes
docker sandbox ls

# Remove all sandboxes (clean slate)
docker sandbox ls | awk 'NR>1 {print $1}' | xargs -I {} docker sandbox rm {}

# Or remove specific sandbox
docker sandbox rm <sandbox-id>

# Kill existing tmux session
tmux kill-session -t loopr-test 2>/dev/null
```

**When to clear sandboxes:**
- Testing authentication flow
- Testing first-time setup/onboarding
- Testing new features that depend on fresh state
- After making changes to sandbox interaction code
- When getting unexpected authentication errors

### Basic tmux Testing Pattern

```bash
# 0. ALWAYS clear sandboxes first (see above)

# 1. Create a new tmux session in your test directory
tmux new-session -d -s loopr-test -c /path/to/test-project

# 2. Send the loopr command
tmux send-keys -t loopr-test "/path/to/loopr plan" Enter

# 3. Wait briefly and capture output (use short sleeps: 0.5-1s is usually enough)
sleep 0.5 && tmux capture-pane -t loopr-test -p

# 4. Interact with prompts
# For Yes/No prompts: send "y" or "n"
tmux send-keys -t loopr-test "y"

# For text input: send text then Enter
tmux send-keys -t loopr-test "5" Enter

# 5. Monitor progress with scrollback
tmux capture-pane -t loopr-test -p -S -100  # Last 100 lines

# 6. Clean up when done
tmux kill-session -t loopr-test
```

### Testing Authentication Flow

```bash
# Start loopr
tmux send-keys -t loopr-test "./loopr plan" Enter

# Wait for auth prompt (0.5s is enough for most prompts)
sleep 0.5 && tmux capture-pane -t loopr-test -p

# Respond to "Authenticate now?" → Yes
tmux send-keys -t loopr-test "y"

# Wait for Claude theme selection
sleep 1 && tmux capture-pane -t loopr-test -p

# Select default theme
tmux send-keys -t loopr-test Enter

# Wait for iterations prompt
sleep 0.5 && tmux capture-pane -t loopr-test -p

# Enter number of iterations
tmux send-keys -t loopr-test "3" Enter

# Monitor the loop
watch -n 2 'tmux capture-pane -t loopr-test -p -S -50'
```

### Sleep Duration Guidelines

- **Initial command**: 0.5s (prompts appear quickly)
- **After keystroke**: 0.5s (UI updates fast)
- **Theme selection**: 1s (Claude initialization takes a moment)
- **During loop**: Check every 5-10s (iterations can take minutes)
- **Never use 2s+ for simple prompts** - it's unnecessarily slow

### Monitoring Long-Running Loops

```bash
# Attach to session to watch live
tmux attach -t loopr-test

# Or poll periodically (detached)
while true; do
  clear
  tmux capture-pane -t loopr-test -p -S -30
  sleep 5
done

# Check git commits being created
watch -n 5 'cd /path/to/test-project && git log --oneline -5'
```

### Common Issues

1. **"Authentication required"** - Remove existing sandbox: `docker sandbox ls` then `docker sandbox rm <id>`
2. **Prompt not responding** - Check if you're sending the right key (`y` vs `Enter`)
3. **Output truncated** - Use `-S -N` flag with larger N value for more scrollback
4. **Session frozen** - Attach and press Ctrl+C to gracefully stop

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
