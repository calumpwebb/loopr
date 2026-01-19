# Loopr: Go + Bubble Tea Implementation Plan

## Overview

Build Loopr as a Go CLI with Bubble Tea for interactive TUI, providing a polished experience for Ralph Playbook loop orchestration.

**Current state**: Python plan exists, starting fresh with Go
**Goal**: Single-binary CLI with interactive prompts and live dashboards

## Why Go + Bubble Tea?

- **Single binary**: No Python/dependencies, just `./loopr`
- **Cross-platform**: Compile for macOS/Linux/Windows
- **Fast startup**: Instant, no interpreter overhead
- **Bubble Tea**: Beautiful interactive prompts and live dashboards
- **Simple deployment**: Drop binary anywhere, it just works

## User Experience

### Command Structure

Only **3 commands**:
- `loopr init` - Interactive setup wizard
- `loopr plan` - Generate implementation plan
- `loopr build` - Build from plan

All options via **interactive prompts**, no CLI flags (except maybe `--help`).

### UX Flows

#### `loopr init` - Setup Wizard

```
┌─────────────────────────────────────────┐
│ Welcome to Loopr!                       │
│ Ralph Loop orchestration for Claude     │
└─────────────────────────────────────────┘

? Where would you like to store loopr files?
  > .loopr

? Which sandbox would you like to use?
  > docker

Checking Docker...
✓ Docker is available

Creating files...
✓ .loopr/PROMPT_build.md
✓ .loopr/PROMPT_plan.md
✓ .loopr/AGENTS.md
✓ .loopr/config.toml
✓ .loopr/specs/README.md
✓ .loopr/specs/example-spec.md
✓ CLAUDE.md

✓ All set! Next: loopr plan
```

**If `.loopr/` already exists:**

```
⚠  .loopr/ directory already exists

? What would you like to do?
  > Overwrite all files
  > Cancel

[If overwrite selected]
⚠  This will overwrite template files!
✓ Files updated
```

**If Docker not available:**

```
✗ Docker is not available

Install Docker Desktop:
  https://www.docker.com/products/docker-desktop

Then run: loopr init
```

#### `loopr plan` / `loopr build` - Live Dashboard

```
Checking sandbox authentication...
✗ Docker sandbox not authenticated

? Authenticate now? (Y/n)
  > Y

Authenticating...
✓ Authenticated!

? How many iterations?
  > 5

┌─────────────────────────────────────────┐
│ Plan Mode - Branch: main                │
│ Iteration: 1/5                          │
└─────────────────────────────────────────┘

[Claude output streams here in real-time...]

Pushing to git...
✓ Pushed to origin/main

====== ITERATION 2/5 ======

[continues...]

┌─────────────────────────────────────────┐
│ ✓ Completed 5/5 iterations              │
│                                         │
│ IMPLEMENTATION_PLAN.md created          │
│ Next: loopr build                       │
└─────────────────────────────────────────┘
```

**Ctrl+C handling:**

```
^C
Interrupted! Exiting gracefully...
✓ Current iteration saved
```

## Architecture

### Project Structure

```
loopr/
├── main.go                      # Entry point
├── go.mod
├── go.sum
├── cmd/
│   ├── init.go                 # Init command (Bubble Tea wizard)
│   ├── plan.go                 # Plan command (auth check + prompt + loop)
│   └── build.go                # Build command (auth check + prompt + loop)
├── internal/
│   ├── config/
│   │   ├── config.go           # Config struct & TOML loading
│   │   └── paths.go            # Path resolution (.loopr/ location)
│   ├── sandbox/
│   │   ├── sandbox.go          # Sandbox interface
│   │   ├── docker.go           # Docker sandbox implementation
│   │   └── auth.go             # Auth validation (haiku test)
│   ├── loop/
│   │   ├── controller.go       # Main loop orchestration
│   │   └── runner.go           # Claude CLI execution
│   ├── git/
│   │   └── operations.go       # Git commands (push, branch detection)
│   ├── ui/
│   │   ├── init_wizard.go      # Bubble Tea init flow
│   │   ├── auth_prompt.go      # Auth confirmation prompt
│   │   ├── iteration_prompt.go # Iteration count input
│   │   ├── dashboard.go        # Live loop dashboard
│   │   └── styles.go           # Lipgloss styles
│   └── templates/
│       └── embed.go            # Embedded template files
└── templates/                   # Template files (embedded at compile time)
    ├── PROMPT_build.md
    ├── PROMPT_plan.md
    ├── AGENTS.md
    ├── config.toml
    ├── CLAUDE.md.template
    └── specs/
        ├── README.md
        └── example-spec.md
```

### Key Components

#### 1. Main Entry (main.go)

Simple CLI dispatcher:

```go
package main

import (
    "fmt"
    "os"
    "loopr/cmd"
)

func main() {
    if len(os.Args) < 2 {
        showHelp()
        os.Exit(0)
    }

    switch os.Args[1] {
    case "init":
        cmd.Init()
    case "plan":
        cmd.Plan()
    case "build":
        cmd.Build()
    case "--help", "-h", "help":
        showHelp()
    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        showHelp()
        os.Exit(1)
    }
}

func showHelp() {
    fmt.Println("Loopr - Ralph Loop orchestration for Claude")
    fmt.Println()
    fmt.Println("Commands:")
    fmt.Println("  loopr init   - Initialize loopr in current directory")
    fmt.Println("  loopr plan   - Generate implementation plan")
    fmt.Println("  loopr build  - Build from plan")
}
```

#### 2. Init Wizard (cmd/init.go + internal/ui/init_wizard.go)

Bubble Tea interactive wizard:

```go
package cmd

import (
    "loopr/internal/ui"
    "loopr/internal/config"
    "loopr/internal/sandbox"
)

func Init() {
    // Check if .loopr exists
    if config.LooprDirExists() {
        // Show overwrite/cancel prompt
        choice := ui.ShowOverwritePrompt()
        if choice == ui.Cancel {
            return
        }
    }

    // Run Bubble Tea wizard
    wizard := ui.NewInitWizard()
    result := wizard.Run() // Returns: looprDir, sandboxType

    // Validate Docker available (fail fast)
    sb := sandbox.NewDocker()
    if !sb.IsAvailable() {
        ui.ShowDockerNotAvailable()
        os.Exit(1)
    }

    // Create files
    if err := createLooprFiles(result.LooprDir, result.SandboxType); err != nil {
        ui.ShowError(err)
        os.Exit(1)
    }

    ui.ShowSuccess()
}
```

**Bubble Tea wizard model:**

```go
package ui

import (
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
)

type InitWizard struct {
    state      int // 0=loopr dir, 1=sandbox type, 2=done
    looprInput textinput.Model
    sandboxList list.Model
    result     InitResult
}

type InitResult struct {
    LooprDir     string
    SandboxType  string
}

func (m InitWizard) Init() tea.Cmd {
    return textinput.Blink
}

func (m InitWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle state transitions, input validation
}

func (m InitWizard) View() string {
    // Render current state
}
```

#### 3. Plan/Build Commands (cmd/plan.go, cmd/build.go)

```go
package cmd

import (
    "loopr/internal/config"
    "loopr/internal/sandbox"
    "loopr/internal/loop"
    "loopr/internal/ui"
)

func Plan() {
    // Load config
    cfg, err := config.Load()
    if err != nil {
        ui.ShowError(err)
        os.Exit(1)
    }

    // Create sandbox
    sb, err := sandbox.New(cfg.Sandbox)
    if err != nil {
        ui.ShowError(err)
        os.Exit(1)
    }

    // Check auth (quick haiku test)
    if !sb.IsAuthenticated() {
        // Show auth prompt
        if ui.PromptAuthenticate() {
            if err := sb.Authenticate(); err != nil {
                ui.ShowError(err)
                os.Exit(1)
            }
            ui.ShowAuthSuccess()
        } else {
            fmt.Println("Authentication required. Run again when ready.")
            os.Exit(0)
        }
    }

    // Prompt for iterations
    maxIterations := ui.PromptIterations()

    // Run loop with live dashboard
    controller := loop.NewController(cfg, sb, "plan", maxIterations)
    controller.Run() // Streams output to terminal
}

func Build() {
    // Same as Plan but mode="build"
}
```

#### 4. Sandbox Interface (internal/sandbox/sandbox.go)

```go
package sandbox

import "os/exec"

type Sandbox interface {
    IsAvailable() bool
    IsAuthenticated() bool
    Authenticate() error
    ExecuteClaude(prompt string, model string) error
}
```

#### 5. Docker Sandbox (internal/sandbox/docker.go)

```go
package sandbox

import (
    "os/exec"
    "os"
)

type DockerSandbox struct{}

func NewDocker() *DockerSandbox {
    return &DockerSandbox{}
}

func (d *DockerSandbox) IsAvailable() bool {
    cmd := exec.Command("docker", "ps")
    err := cmd.Run()
    return err == nil
}

func (d *DockerSandbox) IsAuthenticated() bool {
    // Quick haiku test
    cmd := exec.Command(
        "docker", "sandbox", "run",
        "claude", "-p",
        "--model=claude-haiku-4",
        "--system-prompt=You must reply with only 'OK'. Nothing else.",
        "Say OK",
    )
    err := cmd.Run()
    return err == nil
}

func (d *DockerSandbox) Authenticate() error {
    // Run interactive auth
    cmd := exec.Command(
        "docker", "sandbox", "run",
        "claude", "--version",
    )
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func (d *DockerSandbox) ExecuteClaude(prompt string, model string) error {
    cmd := exec.Command(
        "docker", "sandbox", "run",
        "-w", os.Getwd(),
        "claude", "-p",
        "--dangerously-skip-permissions",
        "--model="+model,
        "--output-format=stream-json",
        prompt,
    )

    // Stream to terminal
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}
```

#### 6. Loop Controller (internal/loop/controller.go)

```go
package loop

import (
    "loopr/internal/config"
    "loopr/internal/sandbox"
    "loopr/internal/git"
)

type Controller struct {
    config    *config.Config
    sandbox   sandbox.Sandbox
    mode      string // "plan" or "build"
    maxIter   int
}

func NewController(cfg *config.Config, sb sandbox.Sandbox, mode string, maxIter int) *Controller {
    return &Controller{
        config:  cfg,
        sandbox: sb,
        mode:    mode,
        maxIter: maxIter,
    }
}

func (c *Controller) Run() error {
    // Load prompt file
    promptFile := c.getPromptFile() // .loopr/PROMPT_plan.md or PROMPT_build.md
    prompt, err := os.ReadFile(promptFile)
    if err != nil {
        return err
    }

    // Get git branch
    branch, err := git.CurrentBranch()
    if err != nil {
        return err
    }

    // Run loop
    for i := 1; i <= c.maxIter; i++ {
        fmt.Printf("\n====== ITERATION %d/%d ======\n\n", i, c.maxIter)

        // Execute Claude
        if err := c.sandbox.ExecuteClaude(string(prompt), "claude-sonnet-4"); err != nil {
            return err
        }

        // Push to git
        fmt.Println("\nPushing to git...")
        if err := git.Push(branch); err != nil {
            return err
        }
        fmt.Println("✓ Pushed to origin/" + branch)
    }

    fmt.Printf("\n✓ Completed %d/%d iterations\n", c.maxIter, c.maxIter)
    return nil
}
```

#### 7. Git Operations (internal/git/operations.go)

```go
package git

import "os/exec"

func CurrentBranch() (string, error) {
    cmd := exec.Command("git", "branch", "--show-current")
    out, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}

func Push(branch string) error {
    cmd := exec.Command("git", "push", "origin", branch)
    err := cmd.Run()
    if err != nil {
        // Try with -u flag (create remote branch)
        cmd = exec.Command("git", "push", "-u", "origin", branch)
        return cmd.Run()
    }
    return nil
}
```

#### 8. UI Components (internal/ui/)

**Styles (internal/ui/styles.go):**

```go
package ui

import "github.com/charmbracelet/lipgloss"

var (
    SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
    ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
    WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
    HeaderStyle  = lipgloss.NewStyle().Bold(true).Border(lipgloss.RoundedBorder())
)
```

**Auth Prompt (internal/ui/auth_prompt.go):**

```go
package ui

import (
    "github.com/charmbracelet/huh"
)

func PromptAuthenticate() bool {
    var authenticate bool

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewConfirm().
                Title("Authenticate now?").
                Value(&authenticate),
        ),
    )

    form.Run()
    return authenticate
}
```

**Iteration Prompt (internal/ui/iteration_prompt.go):**

```go
package ui

import "github.com/charmbracelet/huh"

func PromptIterations() int {
    var iterations int

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("How many iterations?").
                Value(&iterations).
                Validate(func(s string) error {
                    i, err := strconv.Atoi(s)
                    if err != nil || i < 1 {
                        return errors.New("must be >= 1")
                    }
                    return nil
                }),
        ),
    )

    form.Run()
    return iterations
}
```

#### 9. Embedded Templates (internal/templates/embed.go)

```go
package templates

import "embed"

//go:embed *.md *.toml specs/*
var FS embed.FS
```

Files compiled into binary at build time. Extract to `.loopr/` during init.

#### 10. Config Loading (internal/config/config.go)

```go
package config

import (
    "github.com/BurntSushi/toml"
    "path/filepath"
)

type Config struct {
    Sandbox string `toml:"sandbox"`
}

func Load() (*Config, error) {
    var cfg Config

    looprDir := filepath.Join(".", ".loopr")
    configPath := filepath.Join(looprDir, "config.toml")

    if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
        return nil, err
    }

    // Validate
    if cfg.Sandbox != "docker" {
        return nil, errors.New("only 'docker' sandbox supported in v1")
    }

    return &cfg, nil
}

func LooprDirExists() bool {
    info, err := os.Stat(".loopr")
    return err == nil && info.IsDir()
}
```

## Implementation Phases

### Phase 1: Project Setup

```bash
mkdir loopr
cd loopr
go mod init github.com/calumpwebb/loopr
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/huh
go get github.com/BurntSushi/toml
```

Create `main.go` with basic CLI dispatcher.

**Verification**: `go run main.go` shows help

### Phase 2: Template Files

Create `templates/` directory with all Ralph Playbook files:
- PROMPT_build.md
- PROMPT_plan.md
- AGENTS.md
- config.toml
- CLAUDE.md.template
- specs/README.md
- specs/example-spec.md

Create `internal/templates/embed.go` to embed at compile time.

**Verification**: Build and verify templates embedded

### Phase 3: Config & Path Utilities

Implement:
- `internal/config/config.go` - TOML loading and validation
- `internal/config/paths.go` - Path resolution

**Verification**: Unit tests for config loading

### Phase 4: Sandbox Abstraction

Implement:
- `internal/sandbox/sandbox.go` - Interface
- `internal/sandbox/docker.go` - Docker implementation with auth check

**Verification**:
- Test `IsAvailable()` with/without Docker
- Test `IsAuthenticated()` with quick haiku command
- Test `Authenticate()` interactive flow

### Phase 5: Git Operations

Implement `internal/git/operations.go`:
- CurrentBranch()
- Push() with fallback

**Verification**: Unit tests with git repo

### Phase 6: Init Command

Implement:
- `cmd/init.go` - Init command logic
- `internal/ui/init_wizard.go` - Bubble Tea wizard
- File creation from embedded templates

**Verification**: Manual test `loopr init`

### Phase 7: Loop Controller

Implement:
- `internal/loop/controller.go` - Main loop orchestration
- Prompt file loading
- Claude execution
- Git push after each iteration

**Verification**: Integration test with mocked sandbox

### Phase 8: Plan/Build Commands

Implement:
- `cmd/plan.go` - Plan command with auth check + prompts
- `cmd/build.go` - Build command with auth check + prompts
- `internal/ui/auth_prompt.go` - Auth Y/n prompt
- `internal/ui/iteration_prompt.go` - Iteration count input

**Verification**: Manual testing of full flow

### Phase 9: Polish

- Better error messages
- Signal handling (Ctrl+C)
- Help text
- README.md

## Dependencies

```go
require (
    github.com/charmbracelet/bubbletea v0.25.0    // TUI framework
    github.com/charmbracelet/lipgloss v0.9.1      // Styling
    github.com/charmbracelet/huh v0.3.0           // Forms/prompts
    github.com/BurntSushi/toml v1.3.2             // TOML parsing
)
```

## Building & Distribution

```bash
# Build for current platform
go build -o loopr main.go

# Cross-compile for all platforms
GOOS=darwin GOARCH=amd64 go build -o loopr-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o loopr-darwin-arm64
GOOS=linux GOARCH=amd64 go build -o loopr-linux-amd64
GOOS=windows GOARCH=amd64 go build -o loopr-windows-amd64.exe
```

Single binary, no dependencies!

## Critical Implementation Details

### 1. Auth Validation (Quick Haiku Test)

```go
func (d *DockerSandbox) IsAuthenticated() bool {
    cmd := exec.Command(
        "docker", "sandbox", "run",
        "claude", "-p",
        "--model=claude-haiku-4",
        "--system-prompt=You must reply with only 'OK'. Nothing else.",
        "Say OK",
    )

    // Suppress output, just check exit code
    cmd.Stdout = nil
    cmd.Stderr = nil

    err := cmd.Run()
    return err == nil
}
```

Fast, cheap, reliable.

### 2. Streaming Claude Output

```go
func (d *DockerSandbox) ExecuteClaude(prompt string, model string) error {
    cmd := exec.Command(
        "docker", "sandbox", "run",
        "-w", os.Getwd(),
        "claude", "-p",
        "--dangerously-skip-permissions",
        "--model="+model,
        "--output-format=stream-json",
        prompt,
    )

    // CRITICAL: Stream to terminal
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}
```

### 3. Embedded Templates

```go
//go:embed *.md *.toml specs/*
var FS embed.FS

func ExtractTemplates(destDir string) error {
    return fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return os.MkdirAll(filepath.Join(destDir, path), 0755)
        }

        content, err := FS.ReadFile(path)
        if err != nil {
            return err
        }

        return os.WriteFile(filepath.Join(destDir, path), content, 0644)
    })
}
```

### 4. Signal Handling

```go
func (c *Controller) Run() error {
    // Handle Ctrl+C gracefully
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        fmt.Println("\n\nInterrupted! Exiting gracefully...")
        os.Exit(0)
    }()

    // ... loop logic
}
```

## Success Criteria

1. **Single binary** - One file, works anywhere
2. **Fast startup** - Instant response
3. **Interactive UX** - Polished prompts and dashboards
4. **Simple commands** - Only init/plan/build
5. **Fail fast** - Clear errors, no ambiguity
6. **Auth gating** - Automatic auth check with prompts
7. **Streaming output** - Live Claude responses
8. **Cross-platform** - macOS, Linux, Windows

## Future Enhancements

1. **Live dashboard** - Spinner, progress bar during iterations
2. **Color themes** - Light/dark mode support
3. **Config profiles** - Multiple sandbox configs
4. **Status command** - `loopr status` shows current state
5. **Resume capability** - Resume interrupted loops
6. **Additional sandboxes** - E2B, Sprites support

## Notes

- **Go 1.22+** required (for embed improvements)
- **No external deps** at runtime (compiled in)
- **Bubble Tea best practices** - Use huh for simple forms, custom models for complex flows
- **Template updates** - User runs `loopr init` → overwrite to get new templates
- **Minimal UI** - Only ✓ and ✗ symbols, clean output
- **Auth is lazy** - Only check when running plan/build, not during init
