# Loopr

**Ralph Loop orchestration for Claude Code**

Loopr is a single-binary CLI tool that provides a polished, interactive experience for running Ralph Playbook loops with Claude. Built with Go and Bubble Tea, it offers beautiful prompts and live dashboards for planning and building projects.

## Features

- **Single binary**: No dependencies, just download and run
- **Interactive prompts**: Beautiful TUI with Bubble Tea
- **Docker sandbox**: Secure, isolated Claude execution
- **Simple commands**: Only 3 commands to learn
- **Live streaming**: Real-time Claude output
- **Cross-platform**: Works on macOS, Linux, and Windows

## Installation

### Download Binary

Download the latest release for your platform:

```bash
# macOS (Apple Silicon)
curl -L https://github.com/yourusername/loopr/releases/latest/download/loopr-darwin-arm64 -o loopr
chmod +x loopr
sudo mv loopr /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/yourusername/loopr/releases/latest/download/loopr-darwin-amd64 -o loopr
chmod +x loopr
sudo mv loopr /usr/local/bin/

# Linux
curl -L https://github.com/yourusername/loopr/releases/latest/download/loopr-linux-amd64 -o loopr
chmod +x loopr
sudo mv loopr /usr/local/bin/

# Windows
# Download loopr-windows-amd64.exe from releases
```

### Build from Source

Requires Go 1.22 or later:

```bash
git clone https://github.com/yourusername/loopr.git
cd loopr
go build -o loopr main.go
```

## Quick Start

### 1. Initialize

Run the interactive setup wizard:

```bash
loopr init
```

This will:
- Create `.loopr/` directory with configuration
- Set up Docker sandbox
- Install Ralph Playbook templates
- Create example spec files

### 2. Generate Plan

Create an implementation plan:

```bash
loopr plan
```

Loopr will:
- Check Docker authentication (prompts if needed)
- Ask for iteration count
- Run Claude in a loop to generate `IMPLEMENTATION_PLAN.md`
- Automatically push changes to git

### 3. Build

Execute the plan:

```bash
loopr build
```

Loopr will:
- Check Docker authentication
- Ask for iteration count
- Run Claude to implement the plan
- Push changes after each iteration

## How It Works

Loopr orchestrates Ralph Playbook loops using the Docker sandbox for Claude Code:

1. **Setup**: `loopr init` creates configuration and template files
2. **Planning**: `loopr plan` runs Claude iteratively to create a detailed implementation plan
3. **Building**: `loopr build` executes the plan, making incremental changes

Each iteration:
- Loads the prompt file (`.loopr/PROMPT_plan.md` or `.loopr/PROMPT_build.md`)
- Executes Claude with the prompt in Docker sandbox
- Streams output to your terminal in real-time
- Pushes changes to git automatically

## Commands

### `loopr init`

Interactive setup wizard that creates:
- `.loopr/` directory structure
- `PROMPT_plan.md` - Planning prompt template
- `PROMPT_build.md` - Build prompt template
- `AGENTS.md` - Agent behavior instructions
- `config.json` - Sandbox configuration
- `specs/` - Example specification files
- `CLAUDE.md` - Project-specific Claude instructions

Options are selected via interactive prompts (no CLI flags needed).

### `loopr plan`

Generates an implementation plan by running Claude in a loop.

Interactive prompts:
- Docker authentication check (authenticates if needed)
- Number of iterations (how many times to loop)

Output:
- Streams Claude responses in real-time
- Creates `IMPLEMENTATION_PLAN.md`
- Pushes to git after each iteration

### `loopr build`

Executes the implementation plan.

Interactive prompts:
- Docker authentication check (authenticates if needed)
- Number of iterations (how many times to loop)

Output:
- Streams Claude responses in real-time
- Makes incremental code changes
- Pushes to git after each iteration

## Configuration

### `.loopr/config.json`

```json
{
  "sandbox": "docker",
  "model": "claude-sonnet-4"
}
```

Currently only Docker sandbox is supported.

### `.loopr/PROMPT_plan.md`

Template prompt used during planning phase. Customize this to change how Claude generates plans.

### `.loopr/PROMPT_build.md`

Template prompt used during build phase. Customize this to change how Claude implements plans.

### `.loopr/AGENTS.md`

Agent behavior and identity instructions. Defines how Claude should act during loops.

## Requirements

### Docker Desktop

Loopr requires Docker Desktop with Claude Code sandbox support:

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop)
2. Ensure Docker is running
3. Authenticate Claude on first use (Loopr will prompt)

To verify Docker is available:

```bash
docker ps
```

### Claude Code CLI

The Docker sandbox runs the `claude` CLI internally. On first use, Loopr will:
1. Detect if authentication is needed
2. Prompt you to authenticate
3. Run a quick test to verify authentication

No manual setup required.

## Project Structure

After running `loopr init`, your project will have:

```
your-project/
├── .loopr/
│   ├── config.json          # Sandbox configuration
│   ├── PROMPT_plan.md       # Planning prompt template
│   ├── PROMPT_build.md      # Build prompt template
│   ├── AGENTS.md            # Agent behavior instructions
│   └── specs/
│       ├── README.md        # Spec documentation
│       └── example-spec.md  # Example specification
├── CLAUDE.md                # Project-specific Claude instructions
└── [your project files]
```

## Examples

### Basic Workflow

```bash
# Initialize in your project
cd my-project
loopr init

# Create specs describing what to build
vim .loopr/specs/my-feature.md

# Generate implementation plan
loopr plan
# (Select 3 iterations)

# Review IMPLEMENTATION_PLAN.md

# Build the feature
loopr build
# (Select 5 iterations)
```

### Interrupting Loops

Press `Ctrl+C` at any time to gracefully stop:

```
^C
Interrupted! Exiting gracefully...
✓ Current iteration saved
```

Changes are committed and pushed, so you can resume later.

## Troubleshooting

### Docker not available

```
✗ Docker is not available

Install Docker Desktop:
  https://www.docker.com/products/docker-desktop
```

**Solution**: Install Docker Desktop and ensure it's running.

### Docker not authenticated

```
✗ Docker sandbox not authenticated

? Authenticate now? (Y/n)
```

**Solution**: Select `Y` to authenticate. Loopr will run an interactive authentication flow.

### Git push fails

If git push fails, Loopr will try with `-u` flag to create the remote branch.

**Manual fix** if needed:

```bash
git push -u origin <branch-name>
```

## Development

### Building

```bash
go build -o loopr main.go
```

### Testing

```bash
go test ./...
```

### Cross-compiling

```bash
# All platforms
make build-all

# Specific platform
GOOS=darwin GOARCH=arm64 go build -o loopr-darwin-arm64
```

## License

MIT

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms

Inspired by [Ralph Loop](https://github.com/frankbria/ralph-claude-code) for Claude Code.
