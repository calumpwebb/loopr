# Ralph Docker Sandbox Setup

> **Note:** As of v0.2.0, loopr automatically handles Docker sandbox setup. This guide is for manual/advanced setup only. For normal usage, just use `loopr init`, `loopr plan`, and `loopr build`.

## Quick Start

1. **Set your Anthropic API key**:
```bash
export ANTHROPIC_API_KEY='your-api-key-here'
```

2. **Run the sandbox**:
```bash
./ralph-sandbox.sh
```

This will:
- Build a sandboxed Docker container with Claude CLI
- Mount your project at `/workspace`
- Persist Claude credentials between runs
- Drop you into a bash shell

## Inside the Sandbox

Once inside the container, you can:

```bash
# Verify Claude is installed
claude --version

# Login to Claude (first time only)
claude login

# Test Claude
echo "Hello, Claude!" | claude -p

# Run the Ralph loop (once you set it up)
./loop.sh
```

## Security Notes

- The container runs as non-root user `ralph`
- Your project files are mounted but isolated
- Credentials are persisted in a Docker volume
- Network access is NOT restricted by default (add network policies if needed)

## Stop the Sandbox

Just type `exit` to leave the container.

## Clean Up

```bash
# Remove the container
docker-compose -f docker-compose.ralph.yml down

# Remove credentials volume (resets login)
docker volume rm loopr_claude-credentials
```

## Why Use a Sandbox?

The Ralph playbook recommends sandboxing because Ralph needs `--dangerously-skip-permissions`
to run autonomously. Running without a sandbox exposes:
- SSH keys
- AWS/cloud credentials
- Browser cookies
- All files on your machine

**"It's not if it gets popped, it's when. What's the blast radius?"**
