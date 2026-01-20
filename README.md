# Loopr

**Autonomous development orchestration with Claude and Ralph Loop**

Loopr is a single-binary CLI tool that orchestrates autonomous development workflows. It runs Claude in iterative loops to plan and build projects, automatically committing each iteration to git. Built with Go and Bubble Tea, it provides a polished interactive experience for the Ralph Loop methodology.

## Features

- **Single binary**: No dependencies, just download and run
- **Interactive prompts**: Beautiful TUI with Bubble Tea
- **Docker sandbox**: Secure, isolated Claude execution
- **Simple commands**: 6 intuitive commands to learn
- **Task-driven workflow**: Checkbox-based task management
- **Live streaming**: Real-time Claude output
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Embedded prompts**: Consistent behavior, no accidental modifications

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/calumpwebb/loopr/main/scripts/install.sh | bash
```

This will:
- Download the latest release for your platform
- Install to `~/bin` (customizable with `INSTALL_DIR`)
- Suggest PATH updates if needed

### Custom Install Location

```bash
curl -fsSL https://raw.githubusercontent.com/calumpwebb/loopr/main/scripts/install.sh | INSTALL_DIR=/usr/local/bin bash
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
- macOS Intel (AMD64)
- Linux ARM64 (aarch64)
- Linux AMD64 (x86_64)

Coming soon:
- Windows (requires Docker Desktop for Windows testing)

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
- Create `.loopr/` directory with task files
- Create `tasks.md` for tracking work
- Create `context.md` for project info
- Set up archive and PRD directories

### 2. Add Tasks

Edit `.loopr/tasks.md` to add your tasks:

```markdown
# Tasks

- [ ] Implement user authentication (priority: high)
- [ ] Add password reset flow (priority: medium)
- [ ] Create user profile page (priority: low)
```

Or import tasks from a PRD:

```bash
loopr import docs/feature-spec.md
```

### 3. Plan

Analyze and refine your task list:

```bash
loopr plan
```

Loopr will:
- Check Docker authentication (prompts if needed)
- Ask for iteration count
- Analyze your codebase
- Refine and break down tasks in `tasks.md`
- Automatically commit and push changes

### 4. Build

Implement the tasks:

```bash
loopr build
```

Loopr will:
- Check Docker authentication
- Ask for iteration count
- Pick the highest-priority unchecked task
- Implement it completely
- Check off the task and commit
- Push changes after each iteration

### 5. Manage Tasks

```bash
# Check status
loopr status

# Archive completed tasks
loopr archive
```

## How It Works

Loopr orchestrates task-driven development using embedded prompts and the Docker sandbox for Claude Code:

1. **Setup**: `loopr init` creates task management files
2. **Planning**: `loopr plan` runs Claude iteratively to analyze and refine your task list
3. **Building**: `loopr build` picks tasks and implements them one by one
4. **Management**: `loopr status`, `loopr archive`, and `loopr import` help manage your workflow

### Task-Driven Approach

- Tasks are stored in `.loopr/tasks.md` as checkboxes with priorities
- **Plan phase**: Analyzes codebase, refines tasks (no implementation)
- **Build phase**: Implements highest-priority task, checks it off
- All changes are committed and pushed automatically

### Embedded Prompts

- Prompts are baked into the loopr binary (not user-editable files)
- Ensures consistent behavior across projects
- No accidental prompt modifications
- Updates come with new loopr versions

Each iteration:
- Uses embedded prompt based on mode (plan/build)
- Executes Claude with the prompt in Docker sandbox
- Streams output to your terminal in real-time
- Commits and pushes changes automatically

## Commands

### `loopr init`

Interactive setup wizard that creates:
- `.loopr/` directory structure
- `tasks.md` - Task list with checkbox format
- `context.md` - Project context and conventions
- `completed/` - Archive directory for completed tasks
- `prd/` - Directory for PRD files (optional)
- `CLAUDE.md` - Project-specific Claude instructions (at root)

Options are selected via interactive prompts (no CLI flags needed).

### `loopr plan`

Analyzes codebase and refines the task list by running Claude in a loop.

Interactive prompts:
- Docker authentication check (authenticates if needed)
- Number of iterations (how many times to loop)

What it does:
- Reads `.loopr/tasks.md` and `.loopr/context.md`
- Analyzes your codebase
- Refines tasks (breaks down, adds, removes, reprioritizes)
- Updates `.loopr/tasks.md`
- Does NOT implement code (planning only)
- Commits and pushes after each iteration

### `loopr build`

Implements tasks from `.loopr/tasks.md` one at a time.

Interactive prompts:
- Docker authentication check (authenticates if needed)
- Number of iterations (how many times to loop)

What it does:
- Picks highest-priority unchecked task
- Implements the functionality completely
- Runs tests if specified in context.md
- Checks off the task in `.loopr/tasks.md`
- Commits with message referencing the task
- Pushes after each iteration
- Moves to next task in next iteration

### `loopr import <file>`

Imports tasks from a PRD or spec file.

```bash
loopr import docs/feature-spec.md
loopr import .loopr/prd/new-feature.md
```

What it does:
- Reads the specified file
- Extracts actionable tasks
- Assigns priorities (high/medium/low)
- Appends tasks to `.loopr/tasks.md`
- Commits and pushes changes

Use this to:
- Convert PRDs into task lists
- Add work from external documents
- Continuously add new tasks mid-project

### `loopr archive`

Archives completed tasks from `tasks.md` to keep it clean.

```bash
loopr archive
```

What it does:
- Finds all checked tasks `[x]` in `.loopr/tasks.md`
- Moves them to `.loopr/completed/YYYY-MM-DD.md`
- Updates `tasks.md` to remove completed tasks
- Shows summary of archived tasks

Run this when `tasks.md` gets cluttered with completed work.

### `loopr status`

Shows current task status and progress.

```bash
loopr status
```

What it shows:
- Task counts (total, unchecked, completed)
- Progress bar visualization
- Last git commit message
- Git working tree status
- Suggested next steps

Use this to quickly check project progress.

### `loopr update`

Update loopr to the latest version.

```bash
loopr update
```

Checks GitHub releases for updates and installs if available. The update is atomic - if it fails, your current installation remains intact.

## Configuration

### `.loopr/tasks.md`

Your task list in checkbox format:

```markdown
# Tasks

- [ ] Implement user authentication (priority: high)
- [ ] Add password reset flow (priority: medium)
- [x] Setup database schema (priority: high)
- [ ] Create user profile page (priority: low)
```

Format:
- `- [ ]` = Unchecked task (not done)
- `- [x]` = Checked task (completed)
- `(priority: high|medium|low)` = Priority level

The plan phase refines this file, the build phase checks off tasks as it implements them.

### `.loopr/context.md`

Project-specific information for Claude (all sections optional):

```markdown
# Project Context

## Testing
npm test

## Build
npm run build

## Architecture
Next.js 14 with App Router, Supabase for backend

## Conventions
- Use server actions for mutations
- Co-locate tests with source files

## Important Notes
- Auth tokens stored in httpOnly cookies
```

Claude reads this file to understand your project. Update it as your project evolves.

### Model Selection

Loopr uses the **sonnet** model by default for both planning and building. This is hardcoded for simplicity and consistency in v1.

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
│   ├── tasks.md            # Current task list
│   ├── context.md          # Project context
│   ├── completed/          # Archive directory
│   └── prd/                # PRD files (optional)
├── CLAUDE.md               # Project-specific Claude instructions
└── [your project files]
```

As you work:
- `.loopr/tasks.md` gets updated (tasks added, checked off, reprioritized)
- `.loopr/completed/` grows with archived tasks (YYYY-MM-DD.md files)
- Your code changes are committed and pushed automatically

## Examples

### Basic Workflow

```bash
# Initialize in your project
cd my-project
loopr init

# Add tasks manually
vim .loopr/tasks.md
# - [ ] Implement user login (priority: high)
# - [ ] Add logout button (priority: high)
# - [ ] Create user profile page (priority: medium)

# Refine tasks
loopr plan
# (Select 2-3 iterations)

# Review refined .loopr/tasks.md

# Build the features
loopr build
# (Select 5-10 iterations)

# Check progress
loopr status

# Archive completed work
loopr archive
```

### Import from PRD

```bash
# Initialize
loopr init

# Create a PRD
cat > .loopr/prd/authentication.md <<EOF
# Authentication Feature

## Requirements
- Users can register with email/password
- Users can log in
- Users can reset password
- Session management with JWT tokens
EOF

# Import tasks from PRD
loopr import .loopr/prd/authentication.md

# Review extracted tasks
cat .loopr/tasks.md

# Plan and build
loopr plan
loopr build
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
