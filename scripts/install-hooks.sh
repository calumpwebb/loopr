#!/bin/bash
# Install git hooks for the loopr project

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "Installing git hooks..."

# Create pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'HOOK_EOF'
#!/bin/bash
# Pre-commit hook: Format Go code before committing

# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    exit 0
fi

echo "Running go fmt on staged files..."

# Run go fmt on the CLI project
(cd projects/cli && go fmt ./...)

# Add any formatted files back to staging
for file in $STAGED_GO_FILES; do
    if [ -f "$file" ]; then
        git add "$file"
    fi
done

echo "✓ Code formatted successfully"
exit 0
HOOK_EOF

# Make the hook executable
chmod +x "$HOOKS_DIR/pre-commit"

echo "✓ Pre-commit hook installed at $HOOKS_DIR/pre-commit"
echo ""
echo "The hook will:"
echo "  - Run 'go fmt' in projects/cli on staged Go files"
echo "  - Add formatted files back to staging"
echo ""
echo "To bypass the hook (not recommended), use: git commit --no-verify"
