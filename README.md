# Loopr

**Autonomous development orchestration with Claude and Ralph Loop**

Loopr is a single-binary CLI tool that orchestrates autonomous development workflows. It runs Claude in iterative loops to plan and build projects, automatically committing each iteration to git. Built with Go and Bubble Tea, it provides a polished interactive experience for the Ralph Loop methodology.

## Features

- **Single binary**: No dependencies, just download and run
- **Interactive prompts**: Beautiful TUI with Bubble Tea
- **Docker sandbox**: Secure, isolated Claude execution
- **Simple commands**: Only 3 commands to learn
- **Live streaming**: Real-time Claude output
- **Cross-platform**: Works on macOS, Linux, and Windows

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/calumpwebb/loopr/main/install.sh | bash
```

This will:
- Download the latest release for your platform
- Install to `~/bin` (customizable with `INSTALL_DIR`)
- Suggest PATH updates if needed

### Custom Install Location

```bash
curl -fsSL https://raw.githubusercontent.com/calumpwebb/loopr/main/install.sh | INSTALL_DIR=/usr/local/bin bash
```

### Manual Installation

1. Download the latest release for your platform from [GitHub Releases](https://github.com/calumpwebb/loopr/releases)
2. Extract and move to a directory in your PATH:
   ```bash
   mv loopr-darwin-arm64 /usr/local/bin/loopr
   chmod +x /usr/local/bin/loopr
   ```

### Update

```bash
loopr update
```

### Build from Source

Requires Go 1.24+:

```bash
git clone https://github.com/calumpwebb/loopr.git
cd loopr
go build
```

### Supported Platforms

Currently supported:
- macOS Apple Silicon (ARM64)

Coming soon:
- macOS Intel (AMD64)
- Linux ARM64
- Linux AMD64
- Windows

### Check Version

```bash
loopr version
# or
loopr --version
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

### `loopr update`

Update loopr to the latest version.

```bash
loopr update
```

Checks GitHub releases for updates and installs if available. The update is atomic - if it fails, your current installation remains intact.

## Configuration

### `.loopr/config.json`

```json
{
  "sandbox": "docker",
  "looprDir": ".loopr",
  "model": {
    "plan": "sonnet",
    "build": "sonnet"
  }
}
```

Valid model values: `sonnet`, `opus`, `haiku`

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

### Creating a Release

To create a new release, use the release script:

```bash
./scripts/release.sh
```

The script will:
1. Check for uncommitted changes (fails if any exist)
2. Show the current version
3. Prompt for the new version (validates format)
4. Create an annotated git tag
5. Push the tag to GitHub

Once pushed, GitHub Actions will automatically:
- Build the binary for supported platforms
- Create a GitHub release
- Upload the binaries as release assets

**Version format:**
- Release: `v0.1.0`, `v1.2.3`
- Pre-release: `v0.1.0-beta.1`, `v1.0.0-rc.2`

**Permissions:**
Only repository owners/admins can create releases. Collaborators need explicit write access.

## License

MIT

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms

Inspired by [Ralph Loop](https://github.com/frankbria/ralph-claude-code) for Claude Code.
