#!/bin/bash
set -e

# Loopr release script
# Creates and pushes a new version tag

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo "Loopr Release Tool"
echo "=================="
echo ""

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}✗${NC} Uncommitted changes detected"
    echo ""
    echo "Please commit or stash your changes first:"
    echo "  git add ."
    echo "  git commit -m 'your message'"
    echo ""
    exit 1
fi

# Fetch tags
echo "Fetching tags from remote..."
git fetch --tags --quiet

# Get current version
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "none")
echo ""
echo "Current version: $CURRENT_VERSION"
echo ""

# Prompt for new version
read -p "Enter new version (e.g., v0.1.0): " NEW_VERSION

# Validate version format
if ! [[ $NEW_VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    echo ""
    echo -e "${RED}✗${NC} Invalid version format"
    echo ""
    echo "Expected format:"
    echo "  v0.1.0           (release)"
    echo "  v0.1.0-beta.1    (pre-release)"
    echo "  v0.1.0-rc.2      (release candidate)"
    echo ""
    exit 1
fi

# Check if tag already exists
if git rev-parse "$NEW_VERSION" >/dev/null 2>&1; then
    echo ""
    echo -e "${RED}✗${NC} Tag $NEW_VERSION already exists"
    echo ""
    exit 1
fi

# Confirm
echo ""
read -p "Create release $NEW_VERSION? (y/n): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo ""

# Create annotated tag
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
echo -e "${GREEN}✓${NC} Created tag $NEW_VERSION"

# Push tag
git push origin "$NEW_VERSION"
echo -e "${GREEN}✓${NC} Pushed to origin"

echo ""
echo "Release $NEW_VERSION is being built by GitHub Actions"
echo ""
echo "View progress:"
echo "  https://github.com/calumpwebb/loopr/actions"
echo ""
echo "Once complete, the release will be available at:"
echo "  https://github.com/calumpwebb/loopr/releases/tag/$NEW_VERSION"
echo ""
