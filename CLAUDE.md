# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Loopr is a CLI tool for agentic development loops using Claude Code and the Ralph Loop methodology. It runs Claude in iterative loops to plan and build projects, with each iteration automatically committed to git. The tool is written in Go and uses Docker sandbox for secure Claude execution.

**Note**: This is a monorepo. The CLI code lives in `projects/cli/`. All development commands below assume you're in that directory.

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
cd projects/cli
go build -o loopr main.go
```

### Testing
```bash
cd projects/cli

# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests with coverage
go test ./... -cover
```

Current test coverage:
- `internal/config/config_test.go` - Config loading and validation
- `internal/git/operations_test.go` - Git repository operations

### Running Locally
```bash
cd projects/cli
go run main.go <command>
# or after building
./loopr <command>
```

### Cross-Compiling
```bash
cd projects/cli

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

**IMPORTANT**: Always use the release script (NEVER create tags manually). Use tmux for interactive prompts.

```bash
# Start tmux session
tmux new-session -d -s loopr-release -c /Users/calum/Development/loopr/projects/cli

# Run release script
tmux send-keys -t loopr-release "./scripts/release.sh" Enter

# Wait for version prompt
sleep 0.5 && tmux capture-pane -t loopr-release -p

# Enter version (e.g., v0.2.7)
tmux send-keys -t loopr-release "v0.2.7" Enter

# Monitor progress
tmux capture-pane -t loopr-release -p

# Clean up when done
tmux kill-session -t loopr-release
```

The release script will:
1. Check for uncommitted changes
2. Prompt for version (e.g., v0.1.0, v1.0.0-beta.1)
3. Create an annotated git tag
4. Push to GitHub (triggers automated build/release via GitHub Actions)

### Build Provenance Attestations

All release binaries are automatically attested during the GitHub Actions workflow:

**What:** GitHub Artifact Attestations provide cryptographic proof of provenance using Sigstore
**How:** The `actions/attest-build-provenance@v1` action runs after building all binaries
**When:** Attestations are generated before release upload (automatic, no manual steps)
**Compliance:** Meets SLSA Build Level 2 requirements

**For maintainers:**
- Attestations require no manual intervention
- All 4 platform binaries are attested with a single wildcard pattern (`loopr-*`)
- Permissions required: `id-token: write` and `attestations: write` (configured in workflow)
- Verification: Download a release binary and run `gh attestation verify <file> -R calumpwebb/loopr`

**User verification example:**
```bash
# Download a release binary
gh release download v0.2.6 -p "loopr-darwin-arm64"

# Verify attestation
gh attestation verify loopr-darwin-arm64 -R calumpwebb/loopr
```

**Technical details:**
- Signing method: Sigstore ephemeral certificates (keyless signing)
- Attestation format: In-toto provenance (SLSA v1.0)
- Transparency log: Recorded in Sigstore Rekor (immutable audit trail)

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
│   ├── plan.go            # Run planning/build loops (shared runLoop function)
│   ├── import.go          # Import tasks from PRD/spec files
│   ├── archive.go         # Archive completed tasks
│   ├── status.go          # Show task status
│   └── update.go          # Self-update functionality
├── internal/
│   ├── prompts/           # Embedded prompts (plan, build, import)
│   ├── sandbox/           # Docker sandbox abstraction and Claude execution
│   ├── git/               # Git operations (push, branch detection)
│   ├── loop/              # Loop controller (iteration management)
│   └── ui/                # Bubble Tea/Huh prompts and styling
└── templates/             # User-facing templates (tasks.md, context.md, CLAUDE.md)
    └── embed.go           # Go embed directives for templates
```

### User-Facing File Structure (Created by loopr init)

```
.loopr/
├── tasks.md              # Current task list (checkboxes with priorities)
├── context.md           # Project context and conventions
├── completed/           # Archive directory for completed tasks
│   └── YYYY-MM-DD.md   # Tasks archived on specific dates
└── prd/                 # Optional: PRD files to import from
```

### Key Architectural Patterns

#### 1. Command Flow
- `main.go` defines CLI commands using urfave/cli
- Commands delegate to functions in `cmd/` package
- `cmd/plan.go` contains both `Plan()` and `Build()` which call shared `runLoop()` function
- `runLoop()` validates environment, creates sandbox, then delegates to loop controller
- **No config.json needed** - hardcoded to use Docker sandbox and sonnet model

#### 2. Loop Controller Pattern
- `internal/loop/controller.go` orchestrates iterations
- Uses embedded prompts from `internal/prompts` package (not file-based)
- Executes Claude via sandbox for each iteration
- Automatically pushes to git after each iteration
- Handles graceful shutdown on Ctrl+C

#### 3. Embedded Prompts (NEW Architecture)
- Prompts are embedded in Go binary, not user-editable files
- `internal/prompts/plan.go` - Plan prompt (analyzes code, refines tasks)
- `internal/prompts/build.go` - Build prompt (implements tasks)
- `internal/prompts/import.go` - Import prompt (extracts tasks from PRDs)
- **Benefits**: Consistent behavior, no accidental prompt modification, easier updates

#### 4. Task Management
- Tasks stored in `.loopr/tasks.md` as checkboxes with priorities
- Format: `- [ ] Task description (priority: high|medium|low)`
- Plan phase: refines tasks, adds/removes/reprioritizes (does NOT implement)
- Build phase: picks highest-priority task, implements, checks it off
- Archive command: moves completed tasks to `.loopr/completed/YYYY-MM-DD.md`

#### 5. Sandbox Abstraction
- `internal/sandbox/sandbox.go` defines interface
- `internal/sandbox/docker.go` implements Docker-specific logic
- Uses `docker sandbox run` command with working directory mount
- Authentication check via quick haiku test
- Streams stdout/stderr directly to terminal for live output

#### 6. Template System
- Templates embedded at compile time via `//go:embed`
- `templates/embed.go` provides extraction to `.loopr/` directory
- Templates include: tasks.md.template, context.md.template, CLAUDE.md.template
- Init command renames .template files to actual files

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
- Model hardcoded to "sonnet" for consistency (v1)

### Git Operations
- Automatic push after each iteration
- Uses `-u` flag fallback if remote branch doesn't exist
- Requires git repository (checked at start of plan/build)

### Template Extraction
- `loopr init` extracts templates from embedded FS to `.loopr/`
- CLAUDE.md template is extracted to project root (not `.loopr/`)
- Prompts user before overwriting existing `.loopr/` directory

## Code Conventions

### General
- Errors are returned, not logged (except in main.go)
- UI output uses lipgloss styles (defined in `internal/ui/styles.go`)
- All user interaction goes through `internal/ui/prompts.go`
- Git operations isolated in `internal/git/operations.go`
- Sandbox interface allows future non-Docker implementations

### Plan Format
**ALWAYS include a TLDR section at the top of any plan:**

When creating implementation plans (in plan mode or otherwise), structure them with:

```markdown
# TLDR
2-4 sentence summary of what this plan accomplishes and the approach taken.

# [Rest of plan sections...]
```

This helps reviewers quickly understand the scope and strategy before diving into details.

### Commit Messages
**ALWAYS use Conventional Commits format:**

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `style:` - Code style (formatting, missing semicolons, etc)
- `refactor:` - Code change that neither fixes a bug nor adds a feature
- `perf:` - Performance improvement
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks, dependencies, etc
- `ci:` - CI/CD changes

**Breaking changes:** Add `!` after type (e.g., `feat!:`) and describe in footer

**Examples:**
```
feat: add multi-platform support for darwin-amd64, linux-arm64, and linux-amd64
fix: correct version comparison in update command to handle commit hash suffix
docs: update README with new installation instructions
chore: bump dependencies to latest versions
```

**Pre-commit hook:** A git pre-commit hook automatically runs `go fmt` on staged Go files before each commit
