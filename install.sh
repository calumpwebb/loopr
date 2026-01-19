#!/bin/bash
set -e

# Loopr installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/calumpwebb/loopr/main/install.sh | bash
# Custom install dir: curl -fsSL ... | INSTALL_DIR=/custom/path bash

REPO="calumpwebb/loopr"
BINARY_NAME="loopr"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            echo "Supported platforms: macOS (darwin), Linux"
            echo "Coming soon: Windows"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            echo "Supported architectures: amd64, arm64"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
}

# Check if platform is supported
check_platform_supported() {
    # Currently only darwin/arm64 is supported
    if [ "$PLATFORM" != "darwin-arm64" ]; then
        log_error "Platform $PLATFORM is not yet supported"
        echo ""
        echo "Currently supported:"
        echo "  • macOS Apple Silicon (darwin-arm64)"
        echo ""
        echo "Coming soon:"
        echo "  • macOS Intel (darwin-amd64)"
        echo "  • Linux ARM64 (linux-arm64)"
        echo "  • Linux AMD64 (linux-amd64)"
        echo ""
        echo "Follow progress at: https://github.com/$REPO"
        exit 1
    fi
}

# Get latest release version from GitHub API
get_latest_version() {
    local response=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest")
    VERSION=$(echo "$response" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        log_error "Failed to fetch latest version"
        echo "Please check your internet connection or try again later"
        exit 1
    fi

    log_info "Latest version: $VERSION"
}

# Determine install directory
set_install_dir() {
    if [ -n "$INSTALL_DIR" ]; then
        INSTALL_DIR="$INSTALL_DIR"
    else
        INSTALL_DIR="$HOME/bin"
    fi

    # Create directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR"
        log_info "Created directory: $INSTALL_DIR"
    fi
}

# Download and install binary
install_binary() {
    local binary_name="${BINARY_NAME}-${PLATFORM}"
    local download_url="https://github.com/$REPO/releases/download/$VERSION/$binary_name"
    local tmp_file=$(mktemp)

    log_info "Downloading $binary_name..."

    if ! curl -fsSL "$download_url" -o "$tmp_file"; then
        log_error "Download failed"
        echo "URL: $download_url"
        echo "Please check if the release exists or try again later"
        rm -f "$tmp_file"
        exit 1
    fi

    # Make executable
    chmod +x "$tmp_file"

    # Move to install directory
    local install_path="$INSTALL_DIR/$BINARY_NAME"
    if ! mv "$tmp_file" "$install_path"; then
        log_error "Failed to install to $install_path"
        echo "You may need to run with sudo or choose a different directory"
        rm -f "$tmp_file"
        exit 1
    fi

    log_info "Installed to: $install_path"
}

# Check if install dir is in PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_warn "$INSTALL_DIR is not in your PATH"
        echo ""
        echo "Add it to your PATH by adding this to your shell config:"
        echo ""
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""

        # Detect shell and suggest config file
        case "$SHELL" in
            */zsh)
                echo "For zsh, add to: ~/.zshrc"
                ;;
            */bash)
                echo "For bash, add to: ~/.bashrc or ~/.bash_profile"
                ;;
            */fish)
                echo "For fish, run: fish_add_path $INSTALL_DIR"
                ;;
        esac
        echo ""
    else
        log_info "$INSTALL_DIR is in your PATH"
    fi
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local installed_version=$("$BINARY_NAME" version 2>/dev/null || echo "unknown")
        log_info "Installation verified"
        echo ""
        echo "Run '$BINARY_NAME --help' to get started"
    else
        log_warn "Installation complete but '$BINARY_NAME' not found in PATH"
        echo ""
        echo "Add $INSTALL_DIR to your PATH or run directly: $INSTALL_DIR/$BINARY_NAME"
    fi
}

# Main installation flow
main() {
    echo ""
    echo "Loopr Installer"
    echo "==============="
    echo ""

    detect_platform
    log_info "Detected platform: $PLATFORM"

    check_platform_supported
    get_latest_version
    set_install_dir
    install_binary
    check_path
    verify_installation

    echo ""
    log_info "Installation complete!"
    echo ""
}

main
