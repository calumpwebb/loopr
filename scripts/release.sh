#!/bin/bash
set -e

# Loopr release script
# Creates and pushes a new version tag

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Compare two semantic versions
# Returns: 0 if v1 == v2, 1 if v1 > v2, 2 if v1 < v2
compare_versions() {
    local v1=$1
    local v2=$2

    # Strip 'v' prefix and split into base version and pre-release
    v1=${v1#v}
    v2=${v2#v}

    # Split on '-' to separate version from pre-release suffix
    local v1_base=${v1%%-*}
    local v1_pre=${v1#*-}
    [[ "$v1_pre" == "$v1" ]] && v1_pre=""

    local v2_base=${v2%%-*}
    local v2_pre=${v2#*-}
    [[ "$v2_pre" == "$v2" ]] && v2_pre=""

    # Split version numbers
    IFS='.' read -ra v1_parts <<< "$v1_base"
    IFS='.' read -ra v2_parts <<< "$v2_base"

    # Compare major, minor, patch
    for i in 0 1 2; do
        local part1=${v1_parts[$i]:-0}
        local part2=${v2_parts[$i]:-0}

        if (( part1 > part2 )); then
            return 1  # v1 > v2
        elif (( part1 < part2 )); then
            return 2  # v1 < v2
        fi
    done

    # Base versions are equal, check pre-release
    # Release version (no suffix) > pre-release version (has suffix)
    if [[ -z "$v1_pre" && -n "$v2_pre" ]]; then
        return 1  # v1 > v2 (v1 is release, v2 is pre-release)
    elif [[ -n "$v1_pre" && -z "$v2_pre" ]]; then
        return 2  # v1 < v2 (v1 is pre-release, v2 is release)
    elif [[ -n "$v1_pre" && -n "$v2_pre" ]]; then
        # Both have pre-release, compare lexicographically
        if [[ "$v1_pre" > "$v2_pre" ]]; then
            return 1
        elif [[ "$v1_pre" < "$v2_pre" ]]; then
            return 2
        fi
    fi

    return 0  # v1 == v2
}

echo ""
echo "Loopr Release Tool"
echo "=================="
echo ""

# Check if working tree is clean
if [[ -n $(git status --porcelain) ]]; then
    echo -e "${RED}✗${NC} Working tree is not clean"
    echo ""
    git status --short
    echo ""
    echo "Please commit or stash all changes first:"
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

# Compare versions if current version exists
if [[ "$CURRENT_VERSION" != "none" ]]; then
    compare_versions "$NEW_VERSION" "$CURRENT_VERSION"
    result=$?

    if [[ $result -eq 0 ]]; then
        echo ""
        echo -e "${RED}✗${NC} Version $NEW_VERSION is the same as current version $CURRENT_VERSION"
        echo ""
        exit 1
    elif [[ $result -eq 2 ]]; then
        echo ""
        echo -e "${RED}✗${NC} Version $NEW_VERSION is older than current version $CURRENT_VERSION"
        echo ""
        echo "Cannot release an older version. Please use a newer version number."
        echo ""
        exit 1
    fi

    echo -e "${GREEN}✓${NC} Version $NEW_VERSION > $CURRENT_VERSION"
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
